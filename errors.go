package cast

type Error string

func (e Error) String() string {
	return string(e)
}

const (
	tooManyRequests     = Error("too many requests")
	internalServerError = Error("internal server error")
)
