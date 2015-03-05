package memory

import (
	"errors"
	"net/url"

	"../../post"
	storage ".."
)

func init() {
	storage.Register("memory", MemoryBackend{})
}

type MemoryBackend struct{}

func (m MemoryBackend) Open(url *url.URL) (storage.Store, error) {
	store := storage.Store(&MemoryStore{})
	return store, nil
}

type MemoryStore struct{
	posts []post.Post
}

func (s *MemoryStore) FindById(id string) (*post.Post, error) {
	for i, post := range s.posts {
		if post.Id == id {
			return &s.posts[i], nil
		}
	}

	return nil, errors.New("post not found")
}

func (s *MemoryStore) FindAll() ([]post.Post, error) {
	return s.posts, nil
}

func (s *MemoryStore) Create(post post.Post) error {
	s.posts = append(s.posts, post)
	return nil
}

func (s *MemoryStore) Update(updatedPost post.Post) error {
	oldPost, err := s.FindById(updatedPost.Id)
	if err != nil {
		return err
	}

	oldPost.Title = updatedPost.Title
	oldPost.Content = updatedPost.Content
	return nil
}

func (s *MemoryStore) Delete(id string) error {
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

