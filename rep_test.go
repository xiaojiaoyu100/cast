package cast

import (
	"crypto/rand"
	"encoding/xml"
	"net/http"
	"testing"
	"time"
)

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
		Comment string `xml:",comment"`
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

func TestReply_Cost(t *testing.T) {
	cost := 100 * time.Millisecond
	reply := Reply{
		cost: cost,
	}
	if reply.Cost() != cost {
		t.Fatal("Cost() unexpected return.")
	}
}
