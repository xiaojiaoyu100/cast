package cast

// BasicAuth provides info to authenticate
type BasicAuth struct {
	username string
	password string
}

func (ba *BasicAuth) info() (string, string) {
	return ba.username, ba.password
}
