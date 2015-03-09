// a storage whose backend is another instance of gol
package gol

import (
	"errors"
	"net/url"

	storage ".."
	"../../post"
	"../query"
)

type Backend struct{}

type Store struct {
	addr string
}

func init() {
	storage.Register("gol", Backend{})
}

func (b Backend) Open(u *url.URL) (storage.Store, error) {
	store := storage.Store(&Store{u.Host})
	return store, nil
}

func (s *Store) Find(q query.Query) ([]post.Post, error) {
	return nil, errors.New("not implemented")
}

func (s *Store) FindById(id string) (*post.Post, error) {
	// make GET request to s.addr/posts/<id>
	// deserialize response
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
