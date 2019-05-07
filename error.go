package cast

import (
	"io"
	"net"
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
	return err == io.EOF
}

func shouldRetry(err error) bool {
	if isEOF(err) {
		return true
	}
	if isNetworkErr(err) {
		return true
	}
	return false
}
