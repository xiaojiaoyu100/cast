package cast

import (
	"crypto/rand"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
	"testing"
)

func TestWithClient(t *testing.T) {
	client := &http.Client{}
	cast := New(WithClient(client))
	if cast.client != client {
		t.Fatal("fail to initialize cast http client.")
	}
}

func TestWithBasicAuth(t *testing.T) {
	u := make([]byte, 10)
	_, err := rand.Read(u)
	if err != nil {
		t.Fatal(err)
	}
	username := string(u)
	p := make([]byte, 10)
	_, err = rand.Read(p)
	if err != nil {
		t.Fatal(err)
	}
	password := string(p)
	cast := New(WithBasicAuth(username, password))
	if cast.basicAuth.username != username {
		t.Fatal("fail to initialize username.")
	}
	if cast.basicAuth.password != password {
		t.Fatal("fail to initialize password.")
	}
}

func TestWithUrlPrefix(t *testing.T) {
	urlPrefix := "https://www.xiaozhibo.com"
	cast := New(WithUrlPrefix(urlPrefix))
	if cast.urlPrefix != urlPrefix {
		t.Fatal("fail to initialize urlPrefix.")
	}
}

func TestWithHeader(t *testing.T) {
	header := http.Header{}
	header.Set("Accept", "text/plain")
	header.Set("Accept-Charset", "utf-8")
	header.Set("Accept-Encoding", "gzip, deflate")
	cast := New(WithHeader(header))
	if len(cast.header) != len(header) {
		t.Fatal("fail to initialize header.")
	}
}

func TestWithRetryHook(t *testing.T) {
	internalServerErrorHook := func(resp *http.Response) error {
		return nil
	}
	tooManyRequestsHook := func(resp *http.Response) error {
		return nil
	}
	cast := New(WithRetryHook(internalServerErrorHook, tooManyRequestsHook))
	if len(cast.retryHooks) != 2 {
		t.Fatal("fail to initialize retryHooks.")
	}
}

func TestWithRetry(t *testing.T) {
	retry := rand2.Intn(10) + 1
	cast := New(WithRetry(retry))
	if cast.retry != retry {
		t.Fatal("fail to initialize retry.")
	}
}

func TestWithLogger(t *testing.T) {
	logger := log.New(os.Stderr, "", log.Llongfile)
	cast := New(WithLogger(logger))
	if cast.logger != logger {
		t.Fatal("fail to initialize logger.")
	}
}

func TestWithDumpRequestHook(t *testing.T) {
	cast := New(WithDumpRequestHook())
	if cast.dumpRequestHook == nil {
		t.Fatal("fail to initialize dumpRequestHook.")
	}
}

func TestWithDumpResponseHook(t *testing.T) {
	cast := New(WithDumpResponseHook())
	if cast.dumpResponseHook == nil {
		t.Fatal("fail to initialize dumpResponseHook.")
	}
}
