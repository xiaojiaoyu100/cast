package cast

import (
	"log"
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

func WithBearerToken(token string) Setter {
	return func(c *Cast) {
		c.bearerToken = token
	}
}

func WithBaseUrl(u string) Setter {
	return func(c *Cast) {
		c.baseUrl = u
	}
}

func WithHeader(h http.Header) Setter {
	return func(c *Cast) {
		c.header = h
	}
}

func WithRetryHook(hooks ...RetryHook) Setter {
	return func(c *Cast) {
		c.retryHooks = hooks
	}
}

func WithRetry(retry int) Setter {
	return func(c *Cast) {
		c.retry = retry
	}
}

func WithLogger(logger *log.Logger) Setter {
	return func(c *Cast) {
		c.logger = logger
	}
}
