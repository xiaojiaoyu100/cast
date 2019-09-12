package cast

import (
	"encoding/base64"
	"math/rand"
	"testing"
)

func TestBasicAuth_Info(t *testing.T) {
	u := make([]byte, 10)
	p := make([]byte, 10)

	_, err := rand.Read(u)
	ok(t, err)

	username := base64.StdEncoding.EncodeToString(u)

	_, err = rand.Read(p)
	ok(t, err)

	password := base64.StdEncoding.EncodeToString(p)

	basicAuth := new(BasicAuth)
	basicAuth.username = username
	basicAuth.password = password

	a, b := basicAuth.info()

	assert(t, a == username && b == password, "unexpected return")
}
