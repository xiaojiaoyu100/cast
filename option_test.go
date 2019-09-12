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
	cast, err := New(WithBasicAuth(username, password))
	if err != nil {
		t.Fatal("fail to call New()")
	}

	if cast.basicAuth.username != username {
		t.Fatal("fail to initialize username.")
	}
	if cast.basicAuth.password != password {
		t.Fatal("fail to initialize password.")
	}
}

func TestWithBearerToken(t *testing.T) {
	token := "sdsjdsj"
	cast, err := New(WithBearerToken(token))
	if err != nil {
		t.Fatal("fail to call New()")
	}
	assert(t, cast.bearerToken == token, "unexpected token")
}

func TestWithBaseUrl(t *testing.T) {
	u := "https://www.xiaozhibo.com"
	cast, err := New(WithBaseURL(u))
	if err != nil {
		t.Fatal("fail to call New()")
	}
	if cast.baseURL != u {
		t.Fatal("fail to initialize baseUrl.")
	}
}

func TestWithHeader(t *testing.T) {
	header := http.Header{}
	header.Set("Accept", "text/plain")
	header.Set("Accept-Charset", "utf-8")
	header.Set("Accept-Encoding", "gzip, deflate")
	cast, err := New(WithHeader(header))
	if err != nil {
		t.Fatal("fail to call New()")
	}
	if len(cast.header) != len(header) {
		t.Fatal("fail to initialize header.")
	}
}

func TestSetHeader(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal("fail to call New()")
	}
	err = SetHeader("Content-Type", "text/plain", "application/json")(c)
	if err == nil {
		t.Fatal("SetHeader error")
	}
	if len(c.header["Content-Type"]) != 0 {
		t.Fatal("SetHeader error")
	}

	err = SetHeader("Content-Type", "text/plain", "Content-Type", "application/json")(c)
	if err != nil {
		t.Fatal("SetHeader error")
	}
	if len(c.header["Content-Type"]) != 1 {
		t.Fatal("SetHeader error")
	}

	if c.header.Get("Content-Type") != "application/json" {
		t.Fatal("SetHeader error")
	}
}

func TestAddHeader(t *testing.T) {
	c, err := New()
	if err != nil {
		t.Fatal("fail to call New()")
	}
	err = AddHeader("Content-Type", "text/plain", "application/json")(c)
	if err == nil {
		t.Fatal("AddHeader error")
	}
	if len(c.header["Content-Type"]) != 0 {
		t.Fatal("AddHeader error")
	}

	err = AddHeader("Content-Type", "text/plain", "Content-Type", "application/json")(c)
	if err != nil {
		t.Fatal("AddHeader error")
	}
	if len(c.header["Content-Type"]) != 2 {
		t.Fatal("AddHeader error")
	}

	if c.header.Get("Content-Type") != "text/plain" {
		t.Fatal("AddHeader error")
	}
}

func TestWithRetry(t *testing.T) {
	retry := rand2.Intn(10) + 1
	cast, err := New(WithRetry(retry))
	if err != nil {
		t.Fatal("fail to call New()")
	}
	if cast.retry != retry {
		t.Fatal("fail to initialize retry.")
	}
}

func TestWithHttpClientTimeout(t *testing.T) {
	timeout := 1 * time.Second
	cast, err := New(WithHTTPClientTimeout(timeout))
	if err != nil {
		t.Fatal("fail to call New()")
	}
	if cast.httpClientTimeout != timeout {
		t.Fatal("fail to initialize http client timeout")
	}

}
