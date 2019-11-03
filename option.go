package cast

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cep21/circuit/v3"

	"github.com/sirupsen/logrus"
)

// Setter can change the cast instance
type Setter func(cast *Cast) error

// WithBaseURL sets the consistent part of your address.
func WithBaseURL(url string) Setter {
	return func(c *Cast) error {
		c.baseURL = url
		return nil
	}
}

// WithHeader replaces the underlying header.
func WithHeader(h http.Header) Setter {
	return func(c *Cast) error {
		c.header = h
		return nil
	}
}

// SetHeader provides an easy way to set header.
func SetHeader(vv ...string) Setter {
	return func(c *Cast) error {
		if len(vv)%2 != 0 {
			return fmt.Errorf("vv must have even params")
		}
		for i := 0; i < len(vv); i += 2 {
			c.header.Set(vv[i], vv[i+1])
		}
		return nil
	}
}

// AddHeader provides an easy way to add header.
func AddHeader(vv ...string) Setter {
	return func(c *Cast) error {
		if len(vv)%2 != 0 {
			return fmt.Errorf("vv must have even params")
		}
		for i := 0; i < len(vv); i += 2 {
			c.header.Add(vv[i], vv[i+1])
		}
		return nil
	}
}

// WithBasicAuth enables basic auth.
func WithBasicAuth(username, password string) Setter {
	return func(c *Cast) error {
		c.basicAuth = new(BasicAuth)
		c.basicAuth.username = username
		c.basicAuth.password = password
		return nil
	}
}

// WithCookies replaces the underlying cookies which can be sent to server when initiate a request.
func WithCookies(cookies ...*http.Cookie) Setter {
	return func(c *Cast) error {
		c.cookies = cookies
		return nil
	}
}

// WithBearerToken enables bearer authentication.
func WithBearerToken(token string) Setter {
	return func(c *Cast) error {
		c.bearerToken = token
		return nil
	}
}

// WithRetry sets the number of attempts, not counting the normal one.
func WithRetry(retry int) Setter {
	return func(c *Cast) error {
		c.retry = retry
		return nil
	}
}

// WithLinearBackoffStrategy changes the retry strategy called "Linear".
func WithLinearBackoffStrategy(slope time.Duration) Setter {
	return func(c *Cast) error {
		c.stg = linearBackoffStrategy{
			slope: slope,
		}
		return nil
	}
}

// WithConstantBackoffStrategy changes the retry strategy called "Constant".
func WithConstantBackoffStrategy(internal time.Duration) Setter {
	return func(c *Cast) error {
		c.stg = constantBackOffStrategy{
			interval: internal,
		}
		return nil
	}
}

// WithExponentialBackoffStrategy changes the retry strategy called "Exponential".
func WithExponentialBackoffStrategy(base, capacity time.Duration) Setter {
	return func(c *Cast) error {
		c.stg = exponentialBackoffStrategy{
			exponentialBackoff{
				base: base,
				cap:  capacity,
			},
		}
		return nil
	}
}

// WithExponentialBackoffEqualJitterStrategy changes the retry strategy called "Equal Jitter".
func WithExponentialBackoffEqualJitterStrategy(base, capacity time.Duration) Setter {
	return func(c *Cast) error {
		c.stg = exponentialBackoffEqualJitterStrategy{
			exponentialBackoff{
				base: base,
				cap:  capacity,
			},
		}
		return nil
	}
}

// WithExponentialBackoffFullJitterStrategy changes the retry strategy called "Full Jitter".
func WithExponentialBackoffFullJitterStrategy(base, capacity time.Duration) Setter {
	return func(c *Cast) error {
		c.stg = exponentialBackoffFullJitterStrategy{
			exponentialBackoff{
				base: base,
				cap:  capacity,
			},
		}
		return nil
	}
}

// WithExponentialBackoffDecorrelatedJitterStrategy changes the retry strategy called “Decorrelated Jitter”.
func WithExponentialBackoffDecorrelatedJitterStrategy(base, capacity time.Duration) Setter {
	return func(c *Cast) error {
		c.stg = exponentialBackoffDecorrelatedJitterStrategy{
			exponentialBackoff{
				base: base,
				cap:  capacity,
			},
			base,
		}
		return nil
	}
}

// AddRetryHooks adds hooks that can be triggered when in customized conditions
func AddRetryHooks(hooks ...RetryHook) Setter {
	return func(c *Cast) error {
		c.retryHooks = append(c.retryHooks, hooks...)
		return nil
	}
}

// AddResponseHooks adds hooks that can be triggered when a request finished.
func AddResponseHooks(hooks ...responseHook) Setter {
	return func(c *Cast) error {
		c.responseHooks = append(c.responseHooks, hooks...)
		return nil
	}
}

// WithHTTPClientTimeout sets the underlying http client timeout.
func WithHTTPClientTimeout(timeout time.Duration) Setter {
	return func(c *Cast) error {
		c.httpClientTimeout = timeout
		return nil
	}
}

// AddBeforeRequestHook adds a before request hook.
func AddBeforeRequestHook(hks ...BeforeRequestHook) Setter {
	return func(c *Cast) error {
		c.beforeRequestHooks = append(c.beforeRequestHooks, hks...)
		return nil
	}
}

// WithLogHook sets a log callback when condition is achieved.
func WithLogHook(f LogHook) Setter {
	return func(c *Cast) error {
		m := NewMonitor(f)
		c.logger.AddHook(m)
		return nil
	}
}

// WithLogLevel sets log level.
func WithLogLevel(l logrus.Level) Setter {
	return func(c *Cast) error {
		c.logger.SetLevel(l)
		return nil
	}
}

// AddCircuitConfig loads a circuit config.
func AddCircuitConfig(name string, config ...circuit.Config) Setter {
	return func(c *Cast) error {
		if len(config) == 0 {
			config = append(config, circuit.Config{
				Execution: circuit.ExecutionConfig{
					Timeout:               10 * time.Second,
					MaxConcurrentRequests: 1000,
				},
				Fallback: circuit.FallbackConfig{
					MaxConcurrentRequests: 1000,
				},
			})
		}
		_, err := c.h.CreateCircuit(name, config...)
		return err
	}
}

// WithDefaultCircuit sets the default circuit breaker.
func WithDefaultCircuit(name string) Setter {
	return func(c *Cast) error {
		c.defaultCircuitName = name
		return nil
	}
}

// AddRequestHook adds a request hook.
func AddRequestHook(hks ...RequestHook) Setter {
	return func(c *Cast) error {
		c.requestHooks = append(c.requestHooks, hks...)
		return nil
	}
}
