package cast

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"
)

// Response wraps the raw response with attributes.
type Response struct {
	request     *Request
	rawResponse *http.Response
	statusCode  int
	body        []byte
	start       time.Time
	end         time.Time
	cost        time.Duration
	times       int
}

// StatusCode returns http status code.
func (resp *Response) StatusCode() int {
	return resp.statusCode
}

func (resp *Response) Cookies() []*http.Cookie {
	if resp.rawResponse == nil {
		return []*http.Cookie{}
	}
	return resp.rawResponse.Cookies()
}

// Body returns the underlying response body.
func (resp *Response) Body() []byte {
	return resp.body
}

// DecodeFromJson decodes the JSON body into data variable.
func (resp *Response) DecodeFromJson(v interface{}) error {
	return json.Unmarshal(resp.body, &v)
}

// DecodeFromXml decodes the XML body into  data variable.
func (resp *Response) DecodeFromXml(v interface{}) error {
	return xml.Unmarshal(resp.body, &v)
}

// Size returns the length of the body.
func (resp *Response) Size() int64 {
	if resp.rawResponse == nil {
		return 0
	}
	return resp.rawResponse.ContentLength
}

// Header returns the response header.
func (resp *Response) Header() http.Header {
	if resp.rawResponse == nil {
		return http.Header{}
	}
	return resp.rawResponse.Header
}

// StatusOk returns true if http status code is 200, otherwise false.
func (resp *Response) StatusOk() bool {
	return resp.statusCode == http.StatusOK
}

// Start returns the beginning time of a request.
func (resp *Response) Start() time.Time {
	return resp.start
}

// Start returns the end time of a request.
func (resp *Response) End() time.Time {
	return resp.end
}

// Cost returns the cost time of a request.
func (resp *Response) Cost() time.Duration {
	return resp.cost
}
