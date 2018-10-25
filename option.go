package cast

import (
	"net/http"
	"time"
)

type setter func(cast *Cast)

// WithBaseUrl sets the consistent part of your address.
func WithBaseUrl(url string) setter {
	return func(c *Cast) {
		c.baseUrl = url
	}
}

// WithHeader replaces the underlying header.
func WithHeader(h http.Header) setter {
	return func(c *Cast) {
		c.header = h
	}
}

// SetHeader provides an easy way to set header.
func SetHeader(v ...string) setter {
	return func(c *Cast) {
		if len(v)%2 != 0 {
			return
		}
		for i := 0; i < len(v); i++ {
			c.header.Set(v[i], v[i+1])
		}
	}
}

// SetHeader provides an easy way to add header.
func AddHeader(v ...string) setter {
	return func(c *Cast) {
		if len(v)%2 != 0 {
			return
		}
		for i := 0; i < len(v); i++ {
			c.header.Add(v[i], v[i+1])
		}
	}
}

// WithBasicAuth enables basic auth.
func WithBasicAuth(username, password string) setter {
	return func(c *Cast) {
		c.basicAuth = new(BasicAuth)
		c.basicAuth.username = username
		c.basicAuth.password = password
	}
}

// WithCookies replaces the underlying cookies which can be sent to server when initiate a request.
func WithCookies(cookies ...*http.Cookie) setter {
	return func(c *Cast) {
		c.cookies = cookies
	}
}

// WithBearerToken enables bearer authentication.
func WithBearerToken(token string) setter {
	return func(c *Cast) {
		c.bearerToken = token
	}
}

// WithRetry sets the number of attempts, not counting the normal one.
func WithRetry(retry int) setter {
	return func(c *Cast) {
		c.retry = retry
	}
}

// WithLinearBackoffStrategy changes the retry strategy called "Linear".
func WithLinearBackoffStrategy(slope time.Duration) setter {
	return func(c *Cast) {
		c.stg = linearBackoffStrategy{
			slope: slope,
		}
	}
}

// WithConstantBackoffStrategy changes the retry strategy called "Constant".
func WithConstantBackoffStrategy(internal time.Duration) setter {
	return func(c *Cast) {
		c.stg = constantBackOffStrategy{
			interval: internal,
		}
	}
}

// WithConstantBackoffStrategy changes the retry strategy called "Exponential".
func WithExponentialBackoffStrategy(base, cap time.Duration) setter {
	return func(c *Cast) {
		c.stg = exponentialBackoffStrategy{
			exponentialBackoff{
				base: base,
				cap:  cap,
			},
		}
	}
}

// WithExponentialBackoffEqualJitterStrategy changes the retry strategy called "Equal Jitter".
func WithExponentialBackoffEqualJitterStrategy(base, cap time.Duration) setter {
	return func(c *Cast) {
		c.stg = exponentialBackoffEqualJitterStrategy{
			exponentialBackoff{
				base: base,
				cap:  cap,
			},
		}
	}
}

// WithExponentialBackoffFullJitterStrategy changes the retry strategy called "Full Jitter".
func WithExponentialBackoffFullJitterStrategy(base, cap time.Duration) setter {
	return func(c *Cast) {
		c.stg = exponentialBackoffFullJitterStrategy{
			exponentialBackoff{
				base: base,
				cap:  cap,
			},
		}
	}
}

// WithExponentialBackoffFullJitterStrategy changes the retry strategy called “Decorrelated Jitter”.
func WithExponentialBackoffDecorrelatedJitterStrategy(base, cap time.Duration) setter {
	return func(c *Cast) {
		c.stg = exponentialBackoffDecorrelatedJitterStrategy{
			exponentialBackoff{
				base: base,
				cap:  cap,
			},
			base,
		}
	}
}

// AddRetryHooks adds hooks that can be triggered when in customized conditions
func AddRetryHooks(hooks ...RetryHook) setter {
	return func(c *Cast) {
		c.retryHooks = append(c.retryHooks, hooks...)
	}
}

// WithHttpClientTimeout sets the underlying http client timeout.
func WithHttpClientTimeout(timeout time.Duration) setter {
	return func(c *Cast) {
		c.httpClientTimeout = timeout
	}
}
