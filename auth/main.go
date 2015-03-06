package auth

import (
	"errors"
	"fmt"
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

func Open(rawUrl string) (Auth, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	if backend, ok := registeredBackends[u.Scheme]; ok {
		return backend.Open(u)
	} else {
		return nil, errors.New(fmt.Sprint("no such backend:", u.Scheme))
	}
}
