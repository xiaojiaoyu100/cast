package cast

import (
	"crypto/rand"
	"encoding/xml"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestReply_Url(t *testing.T) {
	reply := new(Reply)
	reply.url = "https://google.com"
	assert(t, reply.Url() == "https://google.com", "unexpected Url()")
}

func TestReply_DecodeFromJson(t *testing.T) {
	reply := Reply{
		body: []byte(`{"code": 0, "msg": "ok"}`),
	}
	var temp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := reply.DecodeFromJson(&temp); err != nil {
		t.Fatal(err)
	}
	if temp.Code != 0 && temp.Msg != "ok" {
		t.Fatal("fail to decode json stream.")
	}
}

func TestReply_DecodeFromXml(t *testing.T) {
	reply := new(Reply)
	reply.body = []byte(
		`<person id="13"><name><first>John</first><last>Doe</last></name><age>42</age><Married>false</Married><City>Hanga Roa</City><State>Easter Island</State><!-- Need more details. --></person>`,
	)
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
		Comment   string `xml:",comment"`
	}
	p := Person{}
	err := reply.DecodeFromXml(&p)
	ok(t, err)
	assert(t, p.Id == 13, "unexpected DecodeFromXml")
}

func TestReply_Body(t *testing.T) {
	body := make([]byte, 5)
	_, err := rand.Read(body)
	if err != nil {
		t.Fatal(err)
	}
	reply := Reply{
		body: body,
	}
	if string(reply.Body()) != string(body) {
		t.Fatal("Body() unexpected return")
	}
}

func TestReply_Size(t *testing.T) {
	rep := new(Reply)
	rep.size = 100

	assert(t, rep.Size() == 100, "unexpected Size()")
}

func TestReply_Header(t *testing.T) {
	header := http.Header{
		"Content-Type": []string{"application/json; charset=utf-8"},
	}

	rep := new(Reply)
	rep.header = header

	assert(t, reflect.DeepEqual(rep.Header(), header), "unexpected Header()")
}

func TestReply_StatusOk(t *testing.T) {
	reply := Reply{
		statusCode: http.StatusOK,
	}
	if !reply.StatusOk() {
		t.Fatal("StatusOk() unexpected return.")
	}
}

func TestReply_StatusCode(t *testing.T) {
	reply := Reply{
		statusCode: http.StatusBadRequest,
	}
	if reply.StatusCode() != http.StatusBadRequest {
		t.Fatal("StatusCode() unexpected return.")
	}
}

func TestReply_Start(t *testing.T) {
	start := time.Now().In(time.UTC)
	reply := new(Reply)
	reply.start = start
	assert(t, reply.Start() == start, "unexpected Start()")
}

func TestReply_End(t *testing.T) {
	end := time.Now().In(time.UTC)
	reply := new(Reply)
	reply.end = end
	assert(t, reply.End() == end, "unexpected End()")
}

func TestReply_Cost(t *testing.T) {
	cost := 100 * time.Millisecond
	reply := Reply{
		cost: cost,
	}
	if reply.Cost() != cost {
		t.Fatal("Cost() unexpected return.")
	}
}
