package cast

import (
	"testing"
	"net/http"
	"crypto/rand"
)

func TestReply_DecodeFromJson(t *testing.T) {
	reply := Reply{
		body: []byte(`{"code": 0, "msg": "ok"}`),
	}
	var temp struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := reply.DecodeFromJson(&temp); err != nil {
		t.Fatal(err)
	}
	if temp.Code != 0 && temp.Msg != "ok" {
		t.Fatal("fail to decode json stream.")
	}
}

func TestReply_Body(t *testing.T) {
	body := make([]byte, 5)
	_, err := rand.Read(body)
	if err != nil {
		t.Fatal(err)
	}
	reply := Reply{
		body:        body,
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


