// a storage whose backend is another instance of gol
package gol

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	resp, err := s.doRequest("GET", fmt.Sprintf("/posts/%s", id))
	if err != nil {
		return nil, err
	}

	var p post.Post
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
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

func (s *Store) doRequest(method, path string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", s.addr, path), nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
