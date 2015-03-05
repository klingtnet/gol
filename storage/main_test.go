package storage

import (
	"net/url"
	"testing"
)

type TestStore struct{}

type TestBackend struct{}

func (b TestBackend) Open(u *url.URL) (Store, error) {
	return nil, nil
}

func TestBackendRegistration(t *testing.T) {
	testBackend := TestBackend{}
	Register("test", testBackend)

	backend, ok := registeredBackends["test"]
	if !ok {
		t.Error("backend not registered")
	} else if backend != testBackend {
		t.Error("not the same backend")
	}
}
