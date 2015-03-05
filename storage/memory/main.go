package memory

import (
	"net/url"

	storage ".."
)

func init() {
	storage.Register("memory", MemoryStorage{})
}

type MemoryStorage struct{}

func (m MemoryStorage) Open(url *url.URL) (storage.Store, error) {
	store := MemoryStore{}
	return storage.Store(store), nil
}

type MemoryStore struct{}
