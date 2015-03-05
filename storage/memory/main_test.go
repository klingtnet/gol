package memory

import (
	"net/url"
	"testing"

	"github.com/KLINGTdotNET/gol/post"
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
