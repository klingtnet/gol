package storage

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"../post"
)

var registeredBackends = map[string]Backend{}

func Register(name string, backend Backend) {
	if _, alreadyExists := registeredBackends[name]; !alreadyExists {
		registeredBackends[name] = backend
	} else {
		log.Fatal("duplicate backend", name)
	}
}

func Open(rawUrl string) (Store, error) {
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

type Backend interface {
	Open(url *url.URL) (Store, error)
}

type Store interface {
	FindById(id string) (*post.Post, error)
	FindAll() ([]post.Post, error)

	Create(post post.Post) error
	Update(post post.Post) error
	Delete(id string) error

	// Close() may or may not sync

	// Sync write to disk *now*
}

// Reload = Close + Open
