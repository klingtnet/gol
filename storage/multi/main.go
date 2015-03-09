// a store that writes to multiple backing stores
package multi

import (
	"errors"
	"net/url"

	storage ".."
	"../../post"
	"../query"
)

type Backend struct{}

type Store struct {
	primary     storage.Store
	secondaries []storage.Store
}

func init() {
	storage.Register("multi", Backend{})
}

func (b Backend) Open(u *url.URL) (storage.Store, error) {
	store := storage.Store(&Store{nil, nil})
	return store, nil
}

func (s *Store) Find(q query.Query) ([]post.Post, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) FindById(id string) (*post.Post, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) FindAll() ([]post.Post, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) Create(p post.Post) error {
	return errors.New("not implemented")
}

func (s *Store) Update(p post.Post) error {
	return errors.New("not implemented")
}

func (s *Store) Delete(id string) error {
	return errors.New("not implemented")
}
