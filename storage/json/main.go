package json

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	storage ".."
	"../../post"
	"../memory"
	"../query"
)

type Backend struct{}

type Store struct {
	path          string
	memoryBackend *memory.Store
}

func init() {
	storage.Register("json", Backend{})
}

func (m Backend) Open(u *url.URL) (storage.Store, error) {
	path := u.Host + u.Path

	// initialize with empty array if no such file
	_, err := os.Stat(path)
	if err != nil {
		err = ioutil.WriteFile(path, []byte("[]"), 0644)
		if err != nil {
			return nil, err
		}
	}

	posts, err := readPosts(path)
	if err != nil {
		return nil, err
	}

	store := &Store{
		path:          path,
		memoryBackend: memory.FromPosts(posts),
	}

	return storage.Store(store), nil
}

func readPosts(path string) ([]post.Post, error) {
	postsJson, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var posts []post.Post
	err = json.Unmarshal(postsJson, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func writePosts(path string, posts []post.Post) error {
	post.Reverse(post.ByDate(posts))
	postsJson, err := json.MarshalIndent(posts, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, postsJson, 0644)
}

func (s *Store) Find(q query.Query) ([]post.Post, error) {
	return s.memoryBackend.Find(q)
}

func (s *Store) FindById(id string) (*post.Post, error) {
	return s.memoryBackend.FindById(id)
}

func (s *Store) FindAll() ([]post.Post, error) {
	return s.memoryBackend.FindAll()
}

func (s *Store) Create(post post.Post) error {
	s.memoryBackend.Create(post)
	posts, _ := s.memoryBackend.FindAll()
	return writePosts(s.path, posts)
}

func (s *Store) Update(updatedPost post.Post) error {
	s.memoryBackend.Update(updatedPost)
	posts, _ := s.memoryBackend.FindAll()
	return writePosts(s.path, posts)
}

func (s *Store) Delete(id string) error {
	s.memoryBackend.Delete(id)
	posts, _ := s.memoryBackend.FindAll()
	return writePosts(s.path, posts)
}
