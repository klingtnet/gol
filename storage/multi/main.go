// a store that writes to multiple backing stores
package multi

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	storage ".."
	"../../post"
	"../query"
)

type Backend struct{}

type Store struct {
	primary     storage.Store
	secondaries map[string]storage.Store
}

func init() {
	storage.Register("multi", Backend{})
}

func (b Backend) Open(u *url.URL) (storage.Store, error) {
	primaryUrl := u.Query().Get("primary")
	if primaryUrl == "" {
		return nil, errors.New("no primary store specified")
	}
	primary, err := storage.Open(primaryUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprint("error opening primary store: ", err))
	}

	secondaryUrls := u.Query()["secondary"]
	secondaries := make(map[string]storage.Store, len(secondaryUrls))
	for _, secondaryUrl := range secondaryUrls {
		secondary, err := storage.Open(secondaryUrl)
		if err != nil {
			errors.New(fmt.Sprintf("error opening secondary store '%s': %s", secondaryUrl, err))
		}
		secondaries[secondaryUrl] = secondary
	}

	store := storage.Store(&Store{primary, secondaries})
	return store, nil
}

func (s *Store) Find(q query.Query) ([]post.Post, error) {
	return s.primary.Find(q)
}

func (s *Store) FindById(id string) (*post.Post, error) {
	return s.primary.FindById(id)
}

func (s *Store) FindAll() ([]post.Post, error) {
	return s.primary.FindAll()
}

func (s *Store) Create(p post.Post) error {
	for u, s := range s.secondaries {
		go func(secondary storage.Store) {
			err := secondary.Create(p)
			if err != nil {
				log.Printf("Error: [%s] create: %s", u, err)
			}
		}(s)
	}

	return s.primary.Create(p)
}

func (s *Store) Update(p post.Post) error {
	for u, s := range s.secondaries {
		go func(secondary storage.Store) {
			err := secondary.Update(p)
			if err != nil {
				log.Printf("Error: [%s] update: %s", u, err)
			}
		}(s)
	}

	return s.primary.Update(p)
}

func (s *Store) Delete(id string) error {
	for u, s := range s.secondaries {
		go func(secondary storage.Store) {
			err := secondary.Delete(id)
			if err != nil {
				log.Printf("Error: [%s] delete: %s", u, err)
			}
		}(s)
	}

	return s.primary.Delete(id)
}
