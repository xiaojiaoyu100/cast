package cast

// BasicAuth provides info to authenticate
type BasicAuth struct {
	username string
	password string
}

func (ba *BasicAuth) info() (user, pass string) {
	return ba.username, ba.password
}
