package cast

import (
	"bytes"
	"log"
	"net/http"
	"testing"
)

func TestDumpRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://www.xiaozhibo.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	logger := log.New(&buf, "", log.Llongfile|log.LstdFlags)
	dumpRequest(logger, req)
	t.Log(buf.String())
}
