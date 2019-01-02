package cast

// Error defines cast error
type Error string

func (err Error) Error() string {
	return string(err)
}
