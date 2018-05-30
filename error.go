package cast

type CastError string

func (err CastError) Error() string {
	return string(err)
}
