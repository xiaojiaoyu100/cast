package cast

import (
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func TestBasicAuth_Info(t *testing.T) {
	u := make([]byte, 10)
	p := make([]byte, 10)

	_, err := rand.Read(u)
	if err != nil {
		t.Fatal(err)
	}

	username := base64.StdEncoding.EncodeToString(u)

	_, err = rand.Read(p)
	if err != nil {
		t.Fatal(err)
	}

	password := base64.StdEncoding.EncodeToString(p)

	basicAuth := new(BasicAuth)
	basicAuth.username = username
	basicAuth.password = password

	a, b := basicAuth.info()

	if a != username || b != password {
		t.Fatal("unexpected return")
	}
}
