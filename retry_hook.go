package cast

// RetryHook defines a retry cond.
type RetryHook func(response *Response, err error) bool

var defaultRetryHooks = []RetryHook{
	retry,
}

func retry(_ *Response, err error) bool {
	return shouldRetry(err)
}
