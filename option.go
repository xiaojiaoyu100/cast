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

func WithUrlPrefix(u string) Setter {
	return func(c *Cast) {
		c.urlPrefix = u
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

func WithLogger(log *log.Logger) Setter {
	return func(c *Cast) {
		c.logger = log
	}
}

func WithDumpRequestHook() Setter {
	return func(c *Cast) {
		c.dumpRequestHook = dumpRequest
	}
}

func WithDumpResponseHook() Setter {
	return func(c *Cast) {
		c.dumpResponseHook = dumpResponse
	}
}
