package memory

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"../../post"
)

func TestOpen(t *testing.T) {
	backend := Backend{}
	u, _ := url.Parse("memory://")
	store, _ := backend.Open(u)
	if store == nil {
		t.Fail()
	}
}

func TestCreate(t *testing.T) {
	store := &Store{}
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

func BenchmarkCreate(b *testing.B) {
	store := &Store{}

	for i := 0; i < b.N; i++ {
		store.Create(post.Post{
			Id:      fmt.Sprintf("%d", i),
			Created: time.Now(),
			Title:   "",
			Content: "",
		})
	}
}
