package memory

import (
	"errors"
	"net/url"

	"../../post"
	storage ".."
)

type Backend struct{}

type Store struct{
	posts []post.Post
}

func init() {
	storage.Register("memory", Backend{})
}

func (m Backend) Open(url *url.URL) (storage.Store, error) {
	store := storage.Store(&Store{})
	return store, nil
}

func FromPosts(posts []post.Post) *Store {
	return &Store{
		posts: posts,
	}
}

// `Find` is implemented in `./query.go`

func (s *Store) FindById(id string) (*post.Post, error) {
	for i, post := range s.posts {
		if post.Id == id {
			return &s.posts[i], nil
		}
	}

	return nil, errors.New("post not found")
}

func (s *Store) FindAll() ([]post.Post, error) {
	return s.posts, nil
}

func (s *Store) Create(post post.Post) error {
	s.posts = append(s.posts, post)
	return nil
}

func (s *Store) Update(updatedPost post.Post) error {
	oldPost, err := s.FindById(updatedPost.Id)
	if err != nil {
		return err
	}

	oldPost.Title = updatedPost.Title
	oldPost.Content = updatedPost.Content
	return nil
}

func (s *Store) Delete(id string) error {
	newPosts := make([]post.Post, 0, len(s.posts))
	foundPost := false

	for _, post := range s.posts {
		if post.Id != id {
			newPosts = append(newPosts, post)
		} else {
			foundPost = true
		}
	}

	if !foundPost {
		return errors.New("post not found")
	}

	s.posts = newPosts
	return nil
}
