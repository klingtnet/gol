package nop

import (
	"net/url"

	auth ".."
)

type Backend struct{}

type Auth struct{}

func init() {
	auth.Register("nop", Backend{})
}

func (b Backend) Open(u *url.URL) (auth.Auth, error) {
	return &Auth{}, nil
}

func (a Auth) Login(username, password string) error {
	return nil
}
