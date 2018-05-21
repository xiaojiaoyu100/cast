package cast

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"os"

	"context"

	"github.com/google/go-querystring/query"
	"github.com/jtacoma/uritemplates"
)

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

type Cast struct {
	client     *http.Client
	urlPrefix  string
	api        string
	method     string
	header     http.Header
	queryParam interface{}
	pathParam  map[string]interface{}
	body       reqBody
	basicAuth  *BasicAuth
	retry      int
	strat      backoffStrategy
	retryHooks []RetryHook
	timeout    time.Duration
	logger     *log.Logger
	dumpRequestHook
	dumpResponseHook
}

func New(sl ...Setter) *Cast {
	c := new(Cast)
	c.client = defaultClient
	c.header = make(http.Header)
	c.pathParam = make(map[string]interface{})
	c.logger = log.New(os.Stderr, "CAST ", log.LstdFlags|log.Llongfile)

	for _, s := range sl {
		s(c)
	}

	return c
}

func (c *Cast) WithApi(api string) *Cast {
	c.api = api
	return c
}

func (c *Cast) WithMethod(method string) *Cast {
	c.method = method
	return c
}

func (c *Cast) AppendHeader(header http.Header) *Cast {
	for k, vv := range header {
		for _, v := range vv {
			c.header.Add(k, v)
		}
	}
	return c
}

func (c *Cast) SetHeader(header http.Header) *Cast {
	for k, vv := range header {
		for _, v := range vv {
			c.header.Set(k, v)
		}
	}
	return c
}

func (c *Cast) WithQueryParam(queryParam interface{}) *Cast {
	c.queryParam = queryParam
	return c
}

func (c *Cast) WithPathParam(pathParam map[string]interface{}) *Cast {
	c.pathParam = pathParam
	return c
}

func (c *Cast) WithJsonBody(body interface{}) *Cast {
	c.body = reqJsonBody{
		payload: body,
	}
	return c
}

func (c *Cast) WithUrlEncodedFormBody(body interface{}) *Cast {
	c.body = reqFormUrlEncodedBody{
		payload: body,
	}
	return c
}

func (c *Cast) WithRetry(retry int) *Cast {
	c.retry = retry
	return c
}

func (c *Cast) WithLinearBackoffStrategy(slope time.Duration) *Cast {
	c.strat = linearBackoffStrategy{
		slope: slope,
	}
	return c
}

func (c *Cast) WithConstantBackoffStrategy(internal time.Duration) *Cast {
	c.strat = constantBackOffStrategy{
		interval: internal,
	}
	return c
}

func (c *Cast) WithExponentialBackoffStrategy(base, cap time.Duration) *Cast {
	c.strat = exponentialBackoffStrategy{
		exponentialBackoff{
			base: base,
			cap:  cap,
		},
	}
	return c
}

func (c *Cast) WithExponentialBackoffEqualJitterStrategy(base, cap time.Duration) *Cast {
	c.strat = exponentialBackoffEqualJitterStrategy{
		exponentialBackoff{
			base: base,
			cap:  cap,
		},
	}
	return c
}

func (c *Cast) WithExponentialBackoffFullJitterStrategy(base, cap time.Duration) *Cast {
	c.strat = exponentialBackoffFullJitterStrategy{
		exponentialBackoff{
			base: base,
			cap:  cap,
		},
	}
	return c
}

func (c *Cast) WithExponentialBackoffDecorrelatedJitterStrategy(base, cap time.Duration) *Cast {
	c.strat = exponentialBackoffDecorrelatedJitterStrategy{
		exponentialBackoff{
			base: base,
			cap:  cap,
		},
		base,
	}
	return c
}

func (c *Cast) AddRetryHooks(hooks ...RetryHook) *Cast {
	for _, hook := range hooks {
		c.retryHooks = append(c.retryHooks, hook)
	}
	return c
}

func (c *Cast) WithTimeout(timeout time.Duration) *Cast {
	c.timeout = timeout
	return c
}

func (c *Cast) Request() (*Reply, error) {
	if len(c.pathParam) > 0 {
		tpl, err := uritemplates.Parse(c.api)
		if err != nil {
			c.logger.Printf("ERROR [%v]", err)
			return nil, err
		}
		c.api, err = tpl.Expand(c.pathParam)
		if err != nil {
			c.logger.Printf("ERROR [%v]", err)
			return nil, err
		}
	}

	var (
		reqBody io.Reader
		err     error
	)
	if c.body != nil {
		reqBody, err = c.body.Body()
		if err != nil {
			c.logger.Printf("ERROR [%v]", err)
			return nil, err
		}
		switch c.body.ContentType() {
		case applicaionJson:
			c.SetHeader(http.Header{
				contentType: []string{applicaionJson},
			})
		case formUrlEncoded:
			c.SetHeader(http.Header{
				contentType: []string{formUrlEncoded},
			})
		}
	}
	req, err := http.NewRequest(c.method, c.urlPrefix+c.api, reqBody)
	if err != nil {
		c.logger.Printf("ERROR [%v]", err)
		return nil, err
	}

	for k, vv := range c.header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}

	values, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		c.logger.Printf("ERROR [%v]", err)
		return nil, err
	}

	qValues, err := query.Values(c.queryParam)
	if err != nil {
		c.logger.Printf("ERROR [%v]", err)
		return nil, err
	}
	for k, vv := range qValues {
		for _, v := range vv {
			values.Add(k, v)
		}
	}
	req.URL.RawQuery = values.Encode()

	if c.basicAuth != nil {
		req.SetBasicAuth(c.basicAuth.Info())
	}

	if c.dumpRequestHook != nil {
		c.dumpRequestHook(c.logger, req)
	}

	if c.timeout > 0 {
		ctx, cancel := context.WithCancel(context.TODO())
		_ = time.AfterFunc(c.timeout, func() {
			cancel()
		})
		req = req.WithContext(ctx)
	}

	var (
		resp  *http.Response
		count = 0
	)

	for {

		if count > c.retry {
			break
		}

		resp, err = c.client.Do(req)
		count++

		var isRetry bool
		for _, hook := range c.retryHooks {
			if hook(resp) != nil {
				isRetry = true
				break
			}
		}

		if (isRetry && count <= c.retry+1) || err != nil {
			if resp != nil {
				resp.Body.Close()
			}
			if c.strat != nil {
				<-time.After(c.strat.backoff(count))
				continue
			}
		}

		break
	}

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if c.dumpResponseHook != nil {
		c.dumpResponseHook(c.logger, resp)
	}

	rep := new(Reply)
	rep.statusCode = resp.StatusCode
	repBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Printf("ERROR [%v]", err)
		return nil, err
	}
	rep.body = repBody

	return rep, nil
}
