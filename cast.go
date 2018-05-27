package cast

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var defaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

type Cast struct {
	client             *http.Client
	baseUrl            string
	path               string
	method             string
	header             http.Header
	queryParam         interface{}
	pathParam          map[string]interface{}
	body               reqBody
	basicAuth          *BasicAuth
	bearerToken        string
	retry              int
	strat              backoffStrategy
	timeout            time.Duration
	start              time.Time
	logger             *log.Logger
	retryHooks         []RetryHook
	beforeRequestHooks []BeforeRequestHook
	afterRequestHooks  []AfterRequestHook
	responseHooks      []ResponseHook
	replyHooks         []ReplyHook
}

func New(sl ...Setter) *Cast {
	c := new(Cast)
	c.client = defaultClient
	c.header = make(http.Header)
	c.pathParam = make(map[string]interface{})
	c.logger = log.New(os.Stderr, "[CAST]", log.LstdFlags | log.Llongfile)
	c.beforeRequestHooks = defaultBeforeRequestHooks
	c.afterRequestHooks = defaultAfterRequestHooks

	for _, s := range sl {
		s(c)
	}

	return c
}

func (c *Cast) WithPath(path string) *Cast {
	c.path = path
	return c
}

// Options sets the following http request method to "OPTIONS".
func (c *Cast) Options() *Cast {
	c.method = http.MethodOptions
	return c
}

// Get sets the following http request method to "GET".
func (c *Cast) Get() *Cast {
	c.method = http.MethodGet
	return c
}

// Head sets the following http request method to "HEAD".
func (c *Cast) Head() *Cast {
	c.method = http.MethodHead
	return c
}

// Post sets the following http request method to "POST".
func (c *Cast) Post() *Cast {
	c.method = http.MethodPost
	return c
}

// Put sets the following http request method to "PUT".
func (c *Cast) Put() *Cast {
	c.method = http.MethodPut
	return c
}

// Delete sets the following http request method to "DELETE".
func (c *Cast) Delete() *Cast {
	c.method = http.MethodDelete
	return c
}

// Trace sets the following http request method to "TRACE".
func (c *Cast) Trace() *Cast {
	c.method = http.MethodTrace
	return c
}

// Connect sets the following http request method to "CONNECT".
func (c *Cast) Connect() *Cast {
	c.method = http.MethodConnect
	return c
}

// Patch sets the following http request method to "PATCH".
func (c *Cast) Patch() *Cast {
	c.method = http.MethodPatch
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

func (c *Cast) WithXmlBody(body interface{}) *Cast {
	c.body = reqXmlBody{
		payload: body,
	}
	return c
}

func (c *Cast) WithPlainBody(body string) *Cast {
	c.body = reqPlainBody{
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
	c.retryHooks = append(c.retryHooks, hooks...)
	return c
}

func (c *Cast) AddBeforeRequestHooks(hooks ...BeforeRequestHook) *Cast {
	c.beforeRequestHooks = append(c.beforeRequestHooks, hooks...)
	return c
}

func (c *Cast) AddAfterRequestHooks(hooks ...AfterRequestHook) *Cast {
	c.afterRequestHooks = append(c.afterRequestHooks, hooks...)
	return c
}

func (c *Cast) AddResponseHooks(hooks ...ResponseHook) *Cast {
	c.responseHooks = append(c.responseHooks, hooks...)
	return c
}

func (c *Cast) WithTimeout(timeout time.Duration) *Cast {
	c.timeout = timeout
	return c
}

func (c *Cast) reqBody() (io.Reader, error) {
	var (
		reqBody io.Reader
		err error
	)
	if c.body != nil {
		reqBody, err = c.body.Body()
		if err != nil {
			c.logger.Printf("ERROR [%v]", err)
			return nil, err
		}

	}

	return reqBody, nil
}

func (c *Cast) genReply(request *http.Request) (*Reply, error) {
	var (
		resp  *http.Response
		count = 0
		err error
	)

	for {

		if count > c.retry {
			break
		}

		resp, err = c.client.Do(request)
		count++

		var isRetry bool
		for _, hook := range c.retryHooks {
			if hook(resp) != nil {
				isRetry = true
				break
			}
		}

		if (isRetry && count <= c.retry + 1) || err != nil {
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
		c.logger.Printf("ERROR [%v]", err)
		return nil, err
	}
	defer resp.Body.Close()

	for _, hook := range c.responseHooks {
		if err := hook(c, resp); err != nil {
			return nil, err
		}
	}

	rep := new(Reply)
	rep.url = request.URL.String()
	rep.statusCode = resp.StatusCode
	repBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Printf("ERROR [%v]", err)
		return nil, err
	}
	rep.body = repBody
	rep.size = resp.ContentLength
	rep.header = resp.Header
	rep.cookies = resp.Cookies()
	rep.start = c.start
	rep.end = time.Now().In(time.UTC)
	rep.cost = rep.end.Sub(rep.start)
	rep.times = count

	for _, hook := range c.replyHooks {
		if err := hook(c, rep); err != nil {
			return nil, err
		}
	}

	return rep, nil
}

func (c *Cast) Request() (*Reply, error) {
	reqBody, err := c.reqBody()
	defer putBuffer(reqBody)
	if err != nil {
		return nil, err
	}

	for _, hook := range (c.beforeRequestHooks) {
		if err := hook(c); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(c.method, c.baseUrl + c.path, reqBody)
	if err != nil {
		c.logger.Printf("ERROR [%v]", err)
		return nil, err
	}

	for _, hook := range c.afterRequestHooks {
		req, err = hook(c, req)
		if err != nil {
			return nil, err
		}
	}

	rep, err := c.genReply(req)
	if err != nil {
		return nil, err
	}

	return rep, nil
}
