package cast

import (
	"crypto/rand"
	rand2 "math/rand"
	"net/http"
	"testing"
	"time"
)

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

func TestWithBearerToken(t *testing.T) {
	token := "djsfdeferfrefjnrjfn"
	cast := New(WithBearerToken(token))
	assert(t, cast.bearerToken == token, "unexpected token")
}

func TestWithBaseUrl(t *testing.T) {
	u := "https://www.xiaozhibo.com"
	cast := New(WithBaseUrl(u))
	if cast.baseUrl != u {
		t.Fatal("fail to initialize baseUrl.")
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

func TestWithRetry(t *testing.T) {
	retry := rand2.Intn(10) + 1
	cast := New(WithRetry(retry))
	if cast.retry != retry {
		t.Fatal("fail to initialize retry.")
	}
}

func TestWithHttpClientTimeout(t *testing.T) {
	timeout := 1 * time.Second
	cast := New(WithHttpClientTimeout(timeout))
	if cast.httpClientTimeout != timeout {
		t.Fatal("fail to initialize http client timeout")
	}

}
