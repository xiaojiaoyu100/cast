package cast

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

// Response wraps the raw response with attributes.
type Response struct {
	request     *Request
	rawResponse *http.Response
	statusCode  int
	body        []byte
}

// StatusCode returns http status code.
func (resp *Response) StatusCode() int {
	return resp.statusCode
}

//Cookies returns http cookies.
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

// String returns the underlying body in string.
func (resp *Response) String() string {
	return string(resp.body)
}

// DecodeFromJSON decodes the JSON body into data variable.
func (resp *Response) DecodeFromJSON(v interface{}) error {
	return json.Unmarshal(resp.body, &v)
}

// DecodeFromXML decodes the XML body into  data variable.
func (resp *Response) DecodeFromXML(v interface{}) error {
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

// Success returns true if http status code is in [200,299], otherwise false.
func (resp *Response) Success() bool {
	return resp.statusCode <= 299 && resp.statusCode >= 200
}
