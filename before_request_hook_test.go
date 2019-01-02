package cast

import (
	"encoding/xml"
	"testing"
)

func Test_finalizePathIfAny(t *testing.T) {
	tests := [...]struct {
		path      string
		pathParam map[string]interface{}
		want      string
	}{
		0: {
			path: "/{1}/{2}/{3}",
			pathParam: map[string]interface{}{
				"1": "cd",
				"2": "to",
				"3": "home",
			},
			want: "/cd/to/home",
		},
		1: {
			path: "/{1}/u/{2}",
			pathParam: map[string]interface{}{
				"1": "are",
				"2": "ok",
			},
			want: "/are/u/ok",
		},
	}

	for i, tt := range tests {
		request := new(Request)
		request.pathParam = tt.pathParam
		request.path = tt.path
		err := finalizePathIfAny(nil, request)
		ok(t, err)

		assert(t, request.path == tt.want, "%d: finalizePathIfAny error", i)
	}
}

func Test_setRequestHeader(t *testing.T) {

	c := New()

	jsonBody := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		0,
		"OK",
	}

	plainBody := ""

	xmlBody := struct {
		XMLName xml.Name `xml:"person"`
		Age     int      `xml:"age,attr"`
	}{
		Age: 13,
	}

	tests := [...]struct {
		request *Request
		want    string
	}{
		0: {
			request: c.NewRequest().WithJSONBody(&jsonBody).Post(),
			want:    "application/json; charset=utf-8",
		},
		1: {
			request: c.NewRequest().WithPlainBody(plainBody).Post(),
			want:    "text/plain; charset=utf-8",
		},
		2: {
			request: c.NewRequest().WithXMLBody(xmlBody).Post(),
			want:    "application/xml; charset=utf-8",
		},
	}

	for i, tt := range tests {
		err := setRequestHeader(c, tt.request)
		ok(t, err)
		assert(t, tt.request.header.Get("Content-Type") == tt.want, "%d: unexpected setRequestHeader, got: %s ", i, tt.request.header.Get("Content-Type"))
	}

}
