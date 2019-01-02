package cast

import (
	"net/http"
	"time"
)

// Setter can change the cast instance
type Setter func(cast *Cast)

// WithBaseURL sets the consistent part of your address.
func WithBaseURL(url string) Setter {
	return func(c *Cast) {
		c.baseUURL = url
	}
}

// WithHeader replaces the underlying header.
func WithHeader(h http.Header) Setter {
	return func(c *Cast) {
		c.header = h
	}
}

// SetHeader provides an easy way to set header.
func SetHeader(vv ...string) Setter {
	return func(c *Cast) {
		if len(vv)%2 != 0 {
			return
		}
		for i := 0; i < len(vv); i += 2 {
			c.header.Set(vv[i], vv[i+1])
		}
	}
}

// AddHeader provides an easy way to add header.
func AddHeader(vv ...string) Setter {
	return func(c *Cast) {
		if len(vv)%2 != 0 {
			return
		}
		for i := 0; i < len(vv); i += 2 {
			c.header.Add(vv[i], vv[i+1])
		}
	}
}

// WithBasicAuth enables basic auth.
func WithBasicAuth(username, password string) Setter {
	return func(c *Cast) {
		c.basicAuth = new(BasicAuth)
		c.basicAuth.username = username
		c.basicAuth.password = password
	}
}

// WithCookies replaces the underlying cookies which can be sent to server when initiate a request.
func WithCookies(cookies ...*http.Cookie) Setter {
	return func(c *Cast) {
		c.cookies = cookies
	}
}

// WithBearerToken enables bearer authentication.
func WithBearerToken(token string) Setter {
	return func(c *Cast) {
		c.bearerToken = token
	}
}

// WithRetry sets the number of attempts, not counting the normal one.
func WithRetry(retry int) Setter {
	return func(c *Cast) {
		c.retry = retry
	}
}

// WithLinearBackoffStrategy changes the retry strategy called "Linear".
func WithLinearBackoffStrategy(slope time.Duration) Setter {
	return func(c *Cast) {
		c.stg = linearBackoffStrategy{
			slope: slope,
		}
	}
}

// WithConstantBackoffStrategy changes the retry strategy called "Constant".
func WithConstantBackoffStrategy(internal time.Duration) Setter {
	return func(c *Cast) {
		c.stg = constantBackOffStrategy{
			interval: internal,
		}
	}
}

// WithExponentialBackoffStrategy changes the retry strategy called "Exponential".
func WithExponentialBackoffStrategy(base, cap time.Duration) Setter {
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
func WithExponentialBackoffEqualJitterStrategy(base, cap time.Duration) Setter {
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
func WithExponentialBackoffFullJitterStrategy(base, cap time.Duration) Setter {
	return func(c *Cast) {
		c.stg = exponentialBackoffFullJitterStrategy{
			exponentialBackoff{
				base: base,
				cap:  cap,
			},
		}
	}
}

// WithExponentialBackoffDecorrelatedJitterStrategy changes the retry strategy called “Decorrelated Jitter”.
func WithExponentialBackoffDecorrelatedJitterStrategy(base, cap time.Duration) Setter {
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
func AddRetryHooks(hooks ...RetryHook) Setter {
	return func(c *Cast) {
		c.retryHooks = append(c.retryHooks, hooks...)
	}
}

// WithHTTPClientTimeout sets the underlying http client timeout.
func WithHTTPClientTimeout(timeout time.Duration) Setter {
	return func(c *Cast) {
		c.httpClientTimeout = timeout
	}
}
