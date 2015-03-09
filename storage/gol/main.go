// a storage whose backend is another instance of gol
package gol

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

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
	params := toParams(q)
	resp, err := s.doRequest("GET", fmt.Sprintf("/posts?%s", params.Encode()), nil)
	if err != nil {
		return nil, err
	}

	var posts []post.Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		fmt.Println("decode", err)
		return nil, err
	}

	return posts, nil
}

func toParams(q query.Query) url.Values {
	vals := url.Values{}
	if q.Find != nil {
		vals[q.Find.Name] = []string{q.Find.Value.(string)}
	}
	if q.Start != -1 {
		vals["start"] = []string{fmt.Sprint(q.Start)}
	}
	if q.Count != -1 {
		vals["count"] = []string{fmt.Sprint(q.Count)}
	}
	if len(q.Matches) >= 1 {
		ms := make([]string, 0, len(q.Matches))
		for _, m := range q.Matches {
			ms = append(ms, fmt.Sprintf("%s:%s", m.Name, m.Value.(string)))
		}
		vals["match"] = ms
	}
	if q.RangeStart != nil && q.RangeEnd != nil {
		vals["range"] = []string{fmt.Sprintf("%s,%s", q.RangeStart.Format(time.RFC3339), q.RangeEnd.Format(time.RFC3339))}
	}
	vals["sort"] = []string{q.SortBy}
	vals["reverse"] = []string{fmt.Sprint(q.Reverse)}
	return vals
}

func (s *Store) FindById(id string) (*post.Post, error) {
	resp, err := s.doRequest("GET", fmt.Sprintf("/posts/%s", id), nil)
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
	resp, err := s.doRequest("GET", "/posts", nil)
	if err != nil {
		return nil, err
	}

	var posts []post.Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *Store) Create(p post.Post) error {
	postJson, err := json.Marshal(p)
	if err != nil {
		return err
	}

	resp, err := s.doRequest("POST", "/posts", bytes.NewBuffer(postJson))
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else {
		return errors.New(fmt.Sprintf("unexpected response code: %d (%s)", resp.StatusCode, resp.Status))
	}
}

func (s *Store) Update(p post.Post) error {
	postJson, err := json.Marshal(p)
	if err != nil {
		return err
	}

	resp, err := s.doRequest("POST", fmt.Sprintf("/posts/%s", p.Id), bytes.NewBuffer(postJson))
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else {
		return errors.New(fmt.Sprintf("unexpected response code: %d (%s)", resp.StatusCode, resp.Status))
	}
}

func (s *Store) Delete(id string) error {
	resp, err := s.doRequest("DELETE", fmt.Sprintf("/posts/%s", id), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New(fmt.Sprintf("unexpected response code: %d (%s)", resp.StatusCode, resp.Status))
	}
}

func (s *Store) doRequest(method, path string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s%s", s.addr, path), body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}

	return resp, nil
}
