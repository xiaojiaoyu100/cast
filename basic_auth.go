package cast

type BasicAuth struct {
	username string
	password string
}

func (ba *BasicAuth) Info() (string, string) {
	return ba.username, ba.password
}
