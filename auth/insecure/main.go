package insecure

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"

	auth ".."
)

type Backend struct{}

type Auth struct {
	mapping map[string]string
}

func init() {
	auth.Register("insecure", Backend{})
}

func (b Backend) Open(u *url.URL) (auth.Auth, error) {
	usersJson, err := ioutil.ReadFile(u.Host)
	if err != nil {
		return nil, err
	}

	userCredentials := make(map[string]string)
	err = json.Unmarshal(usersJson, &userCredentials)
	if err != nil {
		return nil, err
	}

	return &Auth{userCredentials}, nil
}

func (a *Auth) Login(username, password string) error {
	if pw, ok := a.mapping[username]; !ok || pw != password {
		return errors.New("invalid credentials")
	}

	return nil
}
