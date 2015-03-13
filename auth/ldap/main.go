package ldap

import (
	"errors"
	"fmt"
	"github.com/vanackere/ldap"
	"net/url"
	"strings"

	auth ".."
)

type Backend struct{}

type Auth struct {
	addr       string
	dnTemplate string
}

func init() {
	auth.Register("ldap", Backend{})
}

func (b Backend) Open(u *url.URL) (auth.Auth, error) {
	dnTemplate := u.Query().Get("dnTemplate")
	if dnTemplate == "" {
		return nil, errors.New("no dnTemplate configured")
	}
	// fix/circumvent uri restrictions
	dnTemplate = strings.Replace(dnTemplate, ":", "=", -1)
	dnTemplate = strings.Replace(dnTemplate, "{}", "%s", -1)

	ldapAuth := Auth{
		addr:       u.Host,
		dnTemplate: dnTemplate,
	}

	return auth.Auth(&ldapAuth), nil
}

func (a Auth) Login(username, password string) error {
	conn, err := ldap.DialTLS("tcp", a.addr, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Bind(fmt.Sprintf(a.dnTemplate, username), password)
	if err != nil {
		return err
	}

	return nil
}
