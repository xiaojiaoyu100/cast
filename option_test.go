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
	cast := New(WithBaseURL(u))
	if cast.baseUURL != u {
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

func TestSetHeader(t *testing.T) {
	c := New()
	SetHeader("Content-Type", "text/plain", "application/json")(c)
	if len(c.header["Content-Type"]) != 0 {
		t.Fatal("SetHeader error")
	}

	SetHeader("Content-Type", "text/plain", "Content-Type", "application/json")(c)
	if len(c.header["Content-Type"]) != 1 {
		t.Fatal("SetHeader error")
	}

	if c.header.Get("Content-Type") != "application/json" {
		t.Fatal("SetHeader error")
	}
}

func TestAddHeader(t *testing.T) {
	c := New()
	AddHeader("Content-Type", "text/plain", "application/json")(c)
	if len(c.header["Content-Type"]) != 0 {
		t.Fatal("AddHeader error")
	}

	AddHeader("Content-Type", "text/plain", "Content-Type", "application/json")(c)
	if len(c.header["Content-Type"]) != 2 {
		t.Fatal("AddHeader error")
	}

	if c.header.Get("Content-Type") != "text/plain" {
		t.Fatal("AddHeader error")
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
	cast := New(WithHTTPClientTimeout(timeout))
	if cast.httpClientTimeout != timeout {
		t.Fatal("fail to initialize http client timeout")
	}

}
