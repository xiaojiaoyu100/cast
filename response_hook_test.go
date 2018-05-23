package cast

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDumpResponse(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://www.xiaozhibo.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(contentType, applicaionJson)
		w.Write([]byte(`{"code": 0, "msg": "ok"}`))
		w.WriteHeader(http.StatusOK)
	}
	handler(rr, req)
	response := rr.Result()
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Llongfile|log.LstdFlags)
	dumpResponse(logger, response)
	t.Log(buf.String())
}
