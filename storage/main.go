package storage

import (
	"errors"
	"fmt"
	"log"
	"net/url"
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
	fmt.Println(registeredBackends, u, u.Scheme)
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
	// FindById(id string) (*Post, error)
	// FindAll() ([]Post, error)

	// Create(post Post) error
	// Update(post Post) error
	// Delete(post Post) error

	// Close() may or may not sync

	// Sync write to disk *now*
}

// Reload = Close + Open
