package cast

import (
	"net/http"
	"time"
)

type profiling struct {
	requestStart      time.Time
	requestDone       time.Time
	requestCost       time.Duration
	dnsStart          time.Time
	dnsDone           time.Time
	dnsCost           time.Duration
	connectStart      time.Time
	connectDone       time.Time
	connectCost       time.Duration
	tlsHandshakeStart time.Time
	tlsHandshakeDone  time.Time
	tlsHandshakeCost  time.Duration
	sendingStart      time.Time
	sendingDone       time.Time
	sendingCost       time.Duration
	waitingStart      time.Time
	waitingDone       time.Time
	waitingCost       time.Duration
	receivingSart     time.Time
	receivingDone     time.Time
	receivingCost     time.Duration
}

// Request is the http.Request wrapper with attributes.
type Request struct {
	path          string
	method        string
	header        http.Header
	queryParam    interface{}
	pathParam     map[string]interface{}
	body          requestBody
	timeout       time.Duration
	remoteAddress string
	prof          profiling
	rawRequest    *http.Request
}

// NewRequest returns an instance of of Request.
func NewRequest() *Request {
	return &Request{
		method:    http.MethodGet,
		header:    make(http.Header),
		pathParam: make(map[string]interface{}),
		prof: profiling{
			requestStart: time.Now().In(time.UTC),
		},
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
func (r *Request) Post() *Request {
	r.method = http.MethodPost
	return r
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

// WithJSONBody creates body with JSON.
func (r *Request) WithJSONBody(body interface{}) *Request {
	r.body = &requestJSONBody{
		payload: body,
	}
	return r
}

// WithXMLBody creates body with XML.
func (r *Request) WithXMLBody(body interface{}) *Request {
	r.body = &requestXMLBody{
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

// WithFormURLEncodedBody creates body with url encoded string.
func (r *Request) WithFormURLEncodedBody(body interface{}) *Request {
	r.body = &requestFormURLEncodedBody{
		payload: body,
	}
	return r
}

// WithMultipartFormDataBody create body with form data
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

// WithHeader replaces the request header.
func (r *Request) WithHeader(header http.Header) *Request {
	r.header = header
	return r
}

// SetHeader sets the key, value pair list.
func (r *Request) SetHeader(vv ...string) *Request {
	if len(vv)%2 != 0 {
		return r
	}

	for i := 0; i < len(vv); i += 2 {
		r.header.Set(vv[i], vv[i+1])
	}
	return r
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

func (r *Request) reqBody() ([]byte, error) {
	if r.body == nil {
		return nil, nil
	}
	body, err := r.body.Body()
	if err != nil {
		return nil, err
	}
	return body, nil
}

// DoHeaderExist whether specified header exists
func (r *Request) DoHeaderExist(h string) bool {
	if r == nil {
		return false
	}
	_, ok := r.header[h]
	return ok
}
