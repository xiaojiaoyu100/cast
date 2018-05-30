package cast

import (
	"io/ioutil"
	"net/http"
	"time"
)

var (
	defaultClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

const (
	defaultDumpBodyLimit int64 = 8192
)

// Cast provides a set of rules to its request.
type Cast struct {
	client             *http.Client
	baseUrl            string
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
	dumpBodyLimit      int64
}

func New(sl ...setter) *Cast {
	c := new(Cast)
	c.client = defaultClient
	c.header = make(http.Header)
	c.beforeRequestHooks = defaultBeforeRequestHooks
	c.requestHooks = defaultRequestHooks
	c.responseHooks = defaultResponseHooks
	c.dumpBodyLimit = defaultDumpBodyLimit

	for _, s := range sl {
		s(c)
	}

	return c
}

// NewRequest returns an instance of Request.
func (c *Cast) NewRequest() *Request {
	return NewRequest()
}

// Do initiates a request.
func (c *Cast) Do(request *Request) (*Response, error) {
	reqBody, err := request.reqBody()
	defer putBuffer(reqBody)
	if err != nil {
		return nil, err
	}

	for _, hook := range c.beforeRequestHooks {
		if err := hook(c, request); err != nil {
			return nil, err
		}
	}

	request.rawRequest, err = http.NewRequest(request.method, c.baseUrl+request.path, reqBody)
	if err != nil {
		globalLogger.printf("ERROR [%v]", err)
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
		globalLogger.printf("ERROR [%v]", err)
		return nil, err
	}
	defer rawResponse.Body.Close()

	resp := new(Response)
	resp.request = request
	resp.rawResponse = rawResponse
	resp.statusCode = rawResponse.StatusCode
	resp.start = request.start
	resp.end = time.Now().In(time.UTC)
	resp.cost = resp.end.Sub(resp.start)
	resp.times = count

	for _, hook := range c.responseHooks {
		if err := hook(c, resp); err != nil {
			return nil, err
		}
	}

	repBody, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		globalLogger.printf("ERROR [%v]", err)
		return nil, err
	}
	resp.body = repBody

	return resp, nil
}
