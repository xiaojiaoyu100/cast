package cast

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
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
	beforeRequestHooks []beforeRequestHook
	requestHooks       []requestHook
	responseHooks      []responseHook
	retryHooks         []RetryHook
	dumpFlag           int
	httpClientTimeout  time.Duration
}

// New returns an instance of Cast
func New(sl ...Setter) *Cast {
	c := new(Cast)
	c.header = make(http.Header)
	c.beforeRequestHooks = defaultBeforeRequestHooks
	c.requestHooks = defaultRequestHooks
	c.responseHooks = defaultResponseHooks
	c.dumpFlag = fStd
	c.httpClientTimeout = 10 * time.Second

	for _, s := range sl {
		s(c)
	}

	c.client = &http.Client{
		Timeout: c.httpClientTimeout,
	}

	roundTripper := http.DefaultTransport
	transport, ok := roundTripper.(*http.Transport)
	if ok {
		transport.MaxIdleConns = 100
		transport.MaxIdleConnsPerHost = 100
	}

	return c
}

// NewRequest returns an instance of Request.
func (c *Cast) NewRequest() *Request {
	return NewRequest()
}

// Do initiates a request.
func (c *Cast) Do(request *Request) (*Response, error) {
	body, err := request.reqBody()
	if err != nil {
		contextLogger.WithError(err)
		return nil, err
	}

	for _, hook := range c.beforeRequestHooks {
		if err := hook(c, request); err != nil {
			return nil, err
		}
	}

	request.rawRequest, err = http.NewRequest(request.method, c.baseURL+request.path, bytes.NewReader(body))
	if err != nil {
		contextLogger.WithError(err)
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

	return rep, nil
}

func (c *Cast) genReply(request *Request) (*Response, error) {
	var (
		rawResponse *http.Response
		count       = 0
		err         error
	)

	for {

		if count > c.retry {
			break
		}

		rawResponse, err = c.client.Do(request.rawRequest)
		count++

		var isRetry bool
		for _, hook := range c.retryHooks {
			if hook(rawResponse) {
				isRetry = true
				break
			}
		}

		if (isRetry && count <= c.retry+1) || err != nil {
			if rawResponse != nil {
				rawResponse.Body.Close()
			}
			if c.stg != nil {
				<-time.After(c.stg.backoff(count))
				continue
			}
		}

		break
	}

	if err != nil {
		contextLogger.WithError(err)
		return nil, err
	}
	defer rawResponse.Body.Close()

	repBody, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		contextLogger.WithError(err)
		return nil, err
	}

	request.prof.requestDone = time.Now().In(time.UTC)
	request.prof.requestCost = request.prof.requestDone.Sub(request.prof.requestStart)

	request.prof.receivingDone = time.Now().In(time.UTC)
	request.prof.receivingCost = request.prof.receivingDone.Sub(request.prof.receivingSart)

	resp := new(Response)
	resp.request = request
	resp.rawResponse = rawResponse
	resp.statusCode = rawResponse.StatusCode
	resp.body = repBody

	for _, hook := range c.responseHooks {
		if err := hook(c, resp); err != nil {
			return nil, err
		}
	}

	return resp, nil
}
