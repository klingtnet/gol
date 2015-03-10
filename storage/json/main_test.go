package json

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"testing"
	"time"

	"../../post"
)

func TestOpen(t *testing.T) {
	backend := Backend{}
	tmpPath, err := ioutil.TempDir("", "gol_json_test")
	if err != nil {
		t.Fatal("Could not create temporary directory", err)
	}

	jsonPath := path.Join(tmpPath, "posts.json")
	u, _ := url.Parse(fmt.Sprintf("json://%s", jsonPath))
	store, err := backend.Open(u)
	if store == nil {
		t.Fatal("could not get store for json backend", err)
	}
}

func BenchmarkCreate(b *testing.B) {
	backend := Backend{}
	tmpPath, err := ioutil.TempDir("", "gol_json_test")
	if err != nil {
		b.Fatal("Could not create temporary directory", err)
	}

	jsonPath := path.Join(tmpPath, "posts.json")
	u, _ := url.Parse(fmt.Sprintf("json://%s", jsonPath))
	store, err := backend.Open(u)
	if store == nil {
		b.Fatal("could not get store for json backend", err)
	}

	for i := 0; i < b.N; i++ {
		store.Create(post.Post{
			Id:      fmt.Sprintf("%d", i),
			Created: time.Now(),
			Title:   "",
			Content: "",
		})
	}
}
