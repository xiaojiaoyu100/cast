package cast

import (
	"net/http"
)

type Setter func(cast *Cast)

func WithClient(client *http.Client) Setter {
	return func(c *Cast) {
		c.client = client
	}
}

func WithBasicAuth(username, password string) Setter {
	return func(c *Cast) {
		c.basicAuth = new(BasicAuth)
		c.basicAuth.username = username
		c.basicAuth.password = password
	}
}




