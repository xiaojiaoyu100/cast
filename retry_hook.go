package cast

import (
	"errors"
	"net/http"
)

type RetryHook func(resp *http.Response) error

func RetryWhenTooManyRequests(resp *http.Response) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		return errors.New(tooManyRequests.String())
	}
	return nil
}

func RetryWhenInternalServerError(resp *http.Response) error {
	if resp.StatusCode == http.StatusInternalServerError {
		return errors.New(internalServerError.String())
	}
	return nil
}
