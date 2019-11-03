package cast

import (
	"io"
	"net"
	"net/url"
)

// Error defines cast error
type Error string

func (err Error) Error() string {
	return string(err)
}

func isNetworkErr(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && (netErr.Temporary() || netErr.Timeout())
}

func isEOF(err error) bool {
	urlErr, ok := err.(*url.Error)
	return ok && urlErr.Err == io.EOF
}

// ShouldRetry returns whether an error needs to be retried.
func ShouldRetry(err error) bool {
	switch {
	case isNetworkErr(err):
		return true
	case isEOF(err):
		return true
	}
	return false
}
