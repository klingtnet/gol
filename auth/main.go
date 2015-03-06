package auth

import (
	"log"
	"net/url"
)

type Backend interface {
	Open(url *url.URL) (Auth, error)
}

type Auth interface {
	Login(username, password string) error
}

var registeredBackends = map[string]Backend{}

func Register(name string, backend Backend) {
	if _, alreadyExists := registeredBackends[name]; !alreadyExists {
		registeredBackends[name] = backend
	} else {
		log.Fatal("duplicate backend:", name)
	}
}
