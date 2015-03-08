package sqlite

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"

	"../../post"
)

func TestOpen(t *testing.T) {
	backend := Backend{}
	// create temporary directory for test sql.db
	tmpPath, err := ioutil.TempDir("", "gol_sqlite_test")
	if err != nil {
		t.Fail()
	}
	u, _ := url.Parse(fmt.Sprintf("sqlite://%s/sqltest.db", tmpPath))
	store, _ := backend.Open(u)
	if store == nil {
		t.Fail()
	}

	err = store.Create(post.Post{Id: "sqlite-test"})
	if err != nil {
		t.Error("could not create post")
		return
	}

	posts, _ := store.FindAll()
	if len(posts) != 1 {
		t.Error("wrong number of posts:", len(posts))
	}

	post, _ := store.FindById("sqlite-test")
	if post == nil {
		t.Error("could not find post")
		return
	}
}

func TestCreate(t *testing.T) {
}
