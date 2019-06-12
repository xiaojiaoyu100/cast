package cast

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cep21/circuit"
	"github.com/cep21/circuit/closers/hystrix"

	"github.com/sirupsen/logrus"
)

const (
	defaultDumpBodyLimit int = 8192
)

// Cast provides a set of rules to its request.
type Cast struct {
	client             *http.Client
	baseURL            string
	header             http.Header
	basicAuth          *BasicAuth
	bearerToken        string
	cookies            []*http.Cookie
	retry              int
	stg                backoffStrategy
	beforeRequestHooks []BeforeRequestHook
	requestHooks       []requestHook
	responseHooks      []responseHook
	retryHooks         []RetryHook
	dumpFlag           int
	httpClientTimeout  time.Duration
	logger             *logrus.Logger
	h                  circuit.Manager
}

// New returns an instance of Cast
func New(sl ...Setter) (*Cast, error) {
	c := new(Cast)
	c.header = make(http.Header)
	c.beforeRequestHooks = defaultBeforeRequestHooks
	c.requestHooks = defaultRequestHooks
	c.responseHooks = defaultResponseHooks
	c.retryHooks = defaultRetryHooks
	c.dumpFlag = fStd
	c.httpClientTimeout = 10 * time.Second
	c.logger = logrus.New()
	c.logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	c.logger.SetReportCaller(true)
	c.logger.SetOutput(os.Stderr)
	c.logger.SetLevel(logrus.WarnLevel)

	configuration := hystrix.Factory{
		ConfigureOpener: hystrix.ConfigureOpener{
			ErrorThresholdPercentage: 60,
			RequestVolumeThreshold:   50,
			RollingDuration:          10 * time.Second,
			Now:                      time.Now,
			NumBuckets:               10,
		},
		ConfigureCloser: hystrix.ConfigureCloser{
			SleepWindow:                  time.Second * 10,
			HalfOpenAttempts:             1,
			RequiredConcurrentSuccessful: 1,
		},
	}

	c.h = circuit.Manager{
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{configuration.Configure},
	}

	for _, s := range sl {
		if err := s(c); err != nil {
			return nil, err
		}
	}

	c.client = &http.Client{
		Timeout: c.httpClientTimeout,
	}

	roundTripper := http.DefaultTransport
	transport, ok := roundTripper.(*http.Transport)
	if ok {
		transport.MaxIdleConns = 1000
		transport.MaxIdleConnsPerHost = 1000
	}

	return c, nil
}

// NewRequest returns an instance of Request.
func (c *Cast) NewRequest() *Request {
	return NewRequest()
}

// Do initiates a request.
func (c *Cast) Do(request *Request) (*Response, error) {
	body, err := request.reqBody()
	if err != nil {
		c.logger.WithError(err).Error("request.reqBody")
		return nil, err
	}

	for _, hook := range c.beforeRequestHooks {
		if err := hook(c, request); err != nil {
			return nil, err
		}
	}

	request.rawRequest, err = http.NewRequest(request.method, c.baseURL+request.path, bytes.NewReader(body))
	if err != nil {
		c.logger.WithError(err).Error("http.NewRequest")
		return nil, err
	}

	for _, hook := range c.requestHooks {
		if err = hook(c, request); err != nil {
			return nil, err
		}
	}

	rep, err := c.genReply(request)
	if err != nil {
		return nil, err
	}

	for _, hook := range c.responseHooks {
		if err := hook(c, rep); err != nil {
			c.logger.WithError(err).Error("hook(c, resp)")
			return nil, err
		}
	}

	return rep, nil
}

func (c *Cast) genReply(request *Request) (*Response, error) {
	var (
		count = 0
		err   error
		resp  *Response
	)
outer:
	for {
		if count > c.retry {
			break outer
		}
		var rawResponse *http.Response

		circuits := c.h.AllCircuits()
		var cb *circuit.Circuit
		switch len(circuits) {
		case 1:
			cb = circuits[0]
		default:
			cb = c.h.GetCircuit(request.configName)
		}
		var fallback bool
		if cb != nil {
			err = cb.Execute(context.TODO(), func(i context.Context) error {
				rawResponse, err = c.client.Do(request.rawRequest)
				if err != nil {
					return err
				}
				return nil
			}, func(i context.Context, e error) error {
				fallback = true
				return e
			})
		} else {
			rawResponse, err = c.client.Do(request.rawRequest)
		}
		if fallback {
			break outer
		}
		count++
		request.prof.requestDone = time.Now().In(time.UTC)
		request.prof.requestCost = request.prof.requestDone.Sub(request.prof.requestStart)
		request.prof.receivingDone = time.Now().In(time.UTC)
		request.prof.receivingCost = request.prof.receivingDone.Sub(request.prof.receivingSart)

		resp = new(Response)
		resp.request = request
		resp.rawResponse = rawResponse
		if rawResponse != nil {
			var repBody []byte
			repBody, err = ioutil.ReadAll(rawResponse.Body)
			if err != nil {
				c.logger.WithError(err).Error("ioutil.ReadAll(rawResponse.Body)")
				return nil, err
			}
			rawResponse.Body.Close()
			resp.body = repBody
			resp.statusCode = rawResponse.StatusCode
		}

		var isRetry bool
		for _, hook := range c.retryHooks {
			if hook(resp, err) {
				isRetry = true
				break
			}
		}

		if isRetry && count < c.retry+1 && c.stg != nil {
			<-time.After(c.stg.backoff(count))
			continue outer
		}

		break outer
	}

	if err != nil {
		c.logger.WithError(err).Error("c.client.Do")
		return nil, err
	}

	return resp, nil
}
