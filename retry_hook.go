package cast

import "net/http"

// RetryHook defines a retry cond.
type RetryHook func(resp *http.Response) bool

// RetryWhenTooManyRequests returns true when http status code is 429, otherwise false.
func RetryWhenTooManyRequests(resp *http.Response) bool {
	if resp.StatusCode == http.StatusTooManyRequests {
		return true
	}
	return false
}

// RetryWhenInternalServerError returns true when http status code is 500, otherwise false.
func RetryWhenInternalServerError(resp *http.Response) bool {
	if resp.StatusCode == http.StatusInternalServerError {
		return true
	}
	return false
}
