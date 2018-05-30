package cast

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestReqJsonBody_ContentType(t *testing.T) {
	reqJsonBody := requestJsonBody{}
	assert(t, reqJsonBody.ContentType() == "application/json; charset=utf-8", "unexpected json ContentType()")
}

func TestReqJsonBody_Body(t *testing.T) {
	type payload struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"
	reqJsonBody := requestJsonBody{
		payload: p,
	}

	body, err := reqJsonBody.Body()
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}

	var b payload
	if err := json.Unmarshal(bytes, &b); err != nil {
		t.Fatal(err)
	}

	if p.Code != b.Code || p.Msg != b.Msg {
		t.Fatal("unexpected body")
	}
}

func TestReqFormUrlEncodedBody_ContentType(t *testing.T) {
	formBody := requestFormUrlEncodedBody{}
	assert(t, formBody.ContentType() == "application/x-www-from-urlencoded", "unexpected form urlencoded ContentType()")
}

func TestReqFormUrlEncodedBody_Body(t *testing.T) {
	type payload struct {
		Code int    `url:"code"`
		Msg  string `url:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"

	req := requestFormUrlEncodedBody{
		payload: p,
	}

	body, err := req.Body()
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		t.Fatal(err)
	}

	values, err := query.Values(&p)
	if err != nil {
		t.Fatal(err)
	}

	if values.Encode() != string(bytes) {
		t.Fatal("unexpected return")
	}
}

func TestReqXmlBody_Body(t *testing.T) {
	type Address struct {
		City, State string
	}
	type Person struct {
		XMLName   xml.Name `xml:"person"`
		Id        int      `xml:"id,attr"`
		FirstName string   `xml:"name>first"`
		LastName  string   `xml:"name>last"`
		Age       int      `xml:"age"`
		Height    float32  `xml:"height,omitempty"`
		Married   bool
		Address
		Comment string `xml:",comment"`
	}
	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	xmlBody := requestXmlBody{
		payload: v,
	}

	body, err := xmlBody.Body()
	ok(t, err)

	bodyBytes, err := ioutil.ReadAll(body)
	ok(t, err)

	var buffer bytes.Buffer
	err = xml.NewEncoder(&buffer).Encode(v)
	ok(t, err)

	t.Log(string(bodyBytes))

	t.Log(string(buffer.String()))

	assert(t, string(bodyBytes) == buffer.String(), "unexpected Body()")

}

func TestReqXmlBody_ContentType(t *testing.T) {
	reqXmlBody := requestXmlBody{}
	assert(t, reqXmlBody.ContentType() == "application/xml; charset=utf-8", "unexpected xml ContentType()")
}

func TestReqPlainBody_Body(t *testing.T) {
	plainBody := requestPlainBody{
		payload: "xssfddfdfdfdfdsfds",
	}

	body, err := plainBody.Body()
	ok(t, err)

	bodyBytes, err := ioutil.ReadAll(body)
	ok(t, err)

	assert(t, string(bodyBytes) == "xssfddfdfdfdfdsfds", "unexpected Body()")
}

func TestReqPlainBody_ContentType(t *testing.T) {
	plainBody := requestPlainBody{}
	assert(t, plainBody.ContentType() == "text/plain; charset=utf-8", "unexpected xml ContentType()")
}
