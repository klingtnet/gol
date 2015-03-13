package insecure

import (
	"errors"
	"net/url"
	"strings"

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
	mapping := make(map[string]string)
	mappings := strings.Split(u.Host, ",")
	for _, m := range mappings {
		s := strings.Split(m, ":")
		if len(s) != 2 {
			return nil, errors.New("invalid user password mapping, format must be 'joe:oops,hey:there'")
		}
		mapping[s[0]] = s[1]
	}

	return &Auth{mapping}, nil
}

func (a *Auth) Login(username, password string) error {
	if pw, ok := a.mapping[username]; !ok || pw != password {
		return errors.New("invalid credentials")
	}

	return nil
}
