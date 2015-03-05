package storage

import (
	"errors"
	"fmt"
	"log"
	"net/url"
)

var storageBackends = map[string]Storage{}

func Register(name string, storageImpl Storage) {
	if _, alreadyExists := storageBackends[name]; !alreadyExists {
		storageBackends[name] = storageImpl
	} else {
		log.Fatal("duplicate backend", name)
	}
}

func Open(rawUrl string) (Store, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	fmt.Println(storageBackends, u, u.Scheme)
	if backend, ok := storageBackends[u.Scheme]; ok {
		return backend.Open(u)
	} else {
		return nil, errors.New(fmt.Sprint("no such backend:", u.Scheme))
	}
}

type Storage interface {
	Open(url *url.URL) (Store, error)
}

type Store interface {
}
