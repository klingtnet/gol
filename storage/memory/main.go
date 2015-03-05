package memory

import (
	"net/url"

	storage ".."
)

func init() {
	storage.Register("memory", MemoryBackend{})
}

type MemoryBackend struct{}

func (m MemoryBackend) Open(url *url.URL) (storage.Store, error) {
	store := MemoryStore{}
	return storage.Store(store), nil
}

type MemoryStore struct{}
