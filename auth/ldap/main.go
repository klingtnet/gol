package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/heyLu/ldap"
	"log"
	"net/url"
	"strings"

	auth ".."
)

type Backend struct{}

type Auth struct {
	addr       string
	insecure   bool
	dnTemplate string
}

func init() {
	auth.Register("ldap", Backend{})
}

func (b Backend) Open(u *url.URL) (auth.Auth, error) {
	insecure := false
	if u.Query().Get("insecure") == "true" {
		log.Println("Warning: insecure ldap connection!")
		insecure = true
	}

	dnTemplate := u.Query().Get("dnTemplate")
	if dnTemplate == "" {
		return nil, errors.New("no dnTemplate configured")
	}
	// fix circumvent uri restrictions
	dnTemplate = strings.Replace(dnTemplate, ",", ";", -1)
	dnTemplate = strings.Replace(dnTemplate, ":", "=", -1)
	dnTemplate = strings.Replace(dnTemplate, "{}", "%s", -1)

	ldapAuth := Auth{
		addr: u.Host,
		dnTemplate: dnTemplate,
		insecure: insecure,
	}

	return auth.Auth(&ldapAuth), nil
}

func (a Auth) Login(username, password string) error {
	var conn *ldap.Conn
	var err error
	if a.insecure {
		tlsConfig := tls.Config{InsecureSkipVerify: true}
		conn, err = ldap.DialSSLWithConfig("tcp", a.addr, &tlsConfig)
	} else {
		conn, err = ldap.DialSSL("tcp", a.addr)
	}
	if err != nil {
		return err
	}

	return conn.Bind(fmt.Sprintf(a.dnTemplate, username), password)
}
