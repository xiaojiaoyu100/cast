package cast

import (
	"io"
	"net/http"
	"time"
)

// Request is the http.Request wrapper with attributes.
type Request struct {
	path       string
	method     string
	header     http.Header
	queryParam interface{}
	pathParam  map[string]interface{}
	body       requestBody
	timeout    time.Duration
	start      time.Time
	rawRequest *http.Request
}

// NewRequest returns an instance of of Request.
func NewRequest() *Request {
	return &Request{
		header:    make(http.Header),
		pathParam: make(map[string]interface{}),
		start:     time.Now().In(time.UTC),
	}
}

// WithPath set the relative or absolute path for the http request
// if the base url don't be provided.
func (r *Request) WithPath(path string) *Request {
	r.path = path
	return r
}

// Options sets the following http request method to "OPTIONS".
func (r *Request) Options() *Request {
	r.method = http.MethodOptions
	return r
}

// Get sets the following http request method to "GET".
func (r *Request) Get() *Request {
	r.method = http.MethodGet
	return r
}

// Head sets the following http request method to "HEAD".
func (r *Request) Head() *Request {
	r.method = http.MethodHead
	return r
}

// Post sets the following http request method to "POST".
func (c *Request) Post() *Request {
	c.method = http.MethodPost
	return c
}

// Put sets the following http request method to "PUT".
func (r *Request) Put() *Request {
	r.method = http.MethodPut
	return r
}

// Delete sets the following http request method to "DELETE".
func (r *Request) Delete() *Request {
	r.method = http.MethodDelete
	return r
}

// Trace sets the following http request method to "TRACE".
func (r *Request) Trace() *Request {
	r.method = http.MethodTrace
	return r
}

// Connect sets the following http request method to "CONNECT".
func (r *Request) Connect() *Request {
	r.method = http.MethodConnect
	return r
}

// Patch sets the following http request method to "PATCH".
func (r *Request) Patch() *Request {
	r.method = http.MethodPatch
	return r
}

// WithQueryParam sets query parameters.
func (r *Request) WithQueryParam(queryParam interface{}) *Request {
	r.queryParam = queryParam
	return r
}

// WithPathParam sets path parameters.
func (r *Request) WithPathParam(pathParam map[string]interface{}) *Request {
	r.pathParam = pathParam
	return r
}

// WithJsonBody creates body with JSON.
func (r *Request) WithJsonBody(body interface{}) *Request {
	r.body = &requestJsonBody{
		payload: body,
	}
	return r
}

// WithXmlBody creates body with XML.
func (r *Request) WithXmlBody(body interface{}) *Request {
	r.body = &requestXmlBody{
		payload: body,
	}
	return r
}

// WithPlainBody creates body with plain text.
func (r *Request) WithPlainBody(body string) *Request {
	r.body = &requestPlainBody{
		payload: body,
	}
	return r
}

// WithFormUrlEncodedBody creates body with url encoded string.
func (r *Request) WithFormUrlEncodedBody(body interface{}) *Request {
	r.body = &requestFormUrlEncodedBody{
		payload: body,
	}
	return r
}

func (r *Request) WithMultipartFormDataBody(formData ...*FormData) *Request {
	r.body = &requestMultipartFormDataBody{
		formData: formData,
	}
	return r
}

// WithTimeout sets the request timeout.
func (r *Request) WithTimeout(timeout time.Duration) *Request {
	r.timeout = timeout
	return r
}

// SetHeader sets the key, value pair list.
func (r *Request) SetHeader(vv ...string) {
	if len(vv)%2 != 0 {
		return
	}

	for i := 0; i < len(vv); i += 2 {
		r.header.Set(vv[i], vv[i+1])
	}
}

// AddHeader adds the key, value pair list.
func (r *Request) AddHeader(vv ...string) {
	if len(vv)%2 != 0 {
		return
	}
	for i := 0; i < len(vv); i += 2 {
		r.header.Add(vv[i], vv[i+1])
	}
}

func (r *Request) reqBody() (io.Reader, error) {
	var (
		reqBody io.Reader
		err     error
	)
	if r.body != nil {
		reqBody, err = r.body.Body()
		if err != nil {
			globalLogger.printf("ERROR [%v]", err)
			return nil, err
		}
	}
	return reqBody, nil
}
