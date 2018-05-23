package cast

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-querystring/query"
)

func TestReqJsonBody_ContentType(t *testing.T) {
	var payload struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	payload.Code = 0
	payload.Msg = "ok"
	reqJsonBody := reqJsonBody{
		payload: payload,
	}
	if reqJsonBody.ContentType() != applicaionJson {
		t.Fatal("unexpected Content-Type")
	}
}

func TestReqJsonBody_Body(t *testing.T) {
	type payload struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"
	reqJsonBody := reqJsonBody{
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
	var payload struct {
		Code int    `url:"code"`
		Msg  string `url:"msg"`
	}
	payload.Code = 0
	payload.Msg = "ok"
	formBody := reqFormUrlEncodedBody{
		payload: payload,
	}
	if formBody.ContentType() != formUrlEncoded {
		t.Fatal("unexpected Content-Type")
	}
}

func TestReqFormUrlEncodedBody_Body(t *testing.T) {
	type payload struct {
		Code int    `url:"code"`
		Msg  string `url:"msg"`
	}

	var p payload
	p.Code = 0
	p.Msg = "ok"

	req := reqFormUrlEncodedBody{
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
