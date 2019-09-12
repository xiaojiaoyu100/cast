package cast

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestReqJsonBody_ContentType(t *testing.T) {
	reqJSONBody := requestJSONBody{}
	assert(t, reqJSONBody.ContentType() == "application/json; charset=utf-8", "unexpected json ContentType()")
}

func TestReqJsonBody_Body(t *testing.T) {
	type payload struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"

	byte, err := json.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}

	req := requestJSONBody{
		payload: byte,
	}

	body, err := req.Body()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))
}

func TestReqFormUrlEncodedBody_ContentType(t *testing.T) {
	formBody := requestFormURLEncodedBody{}
	assert(t, formBody.ContentType() == "application/x-www-form-urlencoded", "unexpected form urlencoded ContentType()")
}

func TestReqFormUrlEncodedBody_Body(t *testing.T) {
	type payload struct {
		Code int    `url:"code"`
		Msg  string `url:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"

	req := requestFormURLEncodedBody{
		payload: p,
	}

	body, err := req.Body()
	if err != nil {
		t.Fatal(err)
	}

	values, err := query.Values(&p)
	if err != nil {
		t.Fatal(err)
	}

	if values.Encode() != string(body) {
		t.Fatal("unexpected return")
	}
}

func TestReqXmlBody_Body(t *testing.T) {
	type Address struct {
		City, State string
	}
	type Person struct {
		XMLName   xml.Name `xml:"person"`
		ID        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	}
	v := &Person{ID: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	xmlBody := requestXMLBody{
		payload: v,
	}

	body, err := xmlBody.Body()
	ok(t, err)

	var buffer bytes.Buffer
	err = xml.NewEncoder(&buffer).Encode(v)
	ok(t, err)

	t.Log(string(buffer.String()))

	assert(t, string(body) == buffer.String(), "unexpected Body()")

}

func TestReqXmlBody_ContentType(t *testing.T) {
	reqXMLBody := requestXMLBody{}
	assert(t, reqXMLBody.ContentType() == "application/xml; charset=utf-8", "unexpected xml ContentType()")
}

func TestReqPlainBody_Body(t *testing.T) {
	plainBody := requestPlainBody{
		payload: "xssfddfdfdfdfdsfds",
	}

	body, err := plainBody.Body()
	ok(t, err)

	assert(t, string(body) == "xssfddfdfdfdfdsfds", "unexpected Body()")
}

func TestReqPlainBody_ContentType(t *testing.T) {
	plainBody := requestPlainBody{}
	assert(t, plainBody.ContentType() == "text/plain; charset=utf-8", "unexpected xml ContentType()")
}
