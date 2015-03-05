package memory

import (
	"net/url"
	"testing"

	"../../post"
)

func TestOpen(t *testing.T) {
	backend := MemoryBackend{}
	u, _ := url.Parse("memory://")
	store, _ := backend.Open(u)
	if store == nil {
		t.Fail()
	}
}

func TestCreate(t *testing.T) {
	store := &MemoryStore{}
	err := store.Create(post.Post{Id: "test"})
	if err != nil {
		t.Error("could not create post")
		return
	}

	posts, _ := store.FindAll()
	if len(posts) != 1 {
		t.Error("wrong number of posts:", len(posts))
	}
}
