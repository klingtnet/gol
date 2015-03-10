package sqlite

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"testing"

	storage ".."
	"../../post"
	tu "../../util/testing"
)

func tSetup(t *testing.T) (storage.Store, func()) {
	backend := Backend{}
	tmpPath, err := ioutil.TempDir("", "gol_sqlite_test")
	if err != nil {
		t.Fatal("could not create temporary directory", err)
	}

	dbPath := path.Join(tmpPath, "sqltest.db")
	u, _ := url.Parse(fmt.Sprintf("sqlite://%s", dbPath))
	store, err := backend.Open(u)
	if store == nil {
		t.Fatal("could not get store for sqlite backend", err)
	}

	return store, func() {
		store.Close()
		os.RemoveAll(tmpPath)
	}
}

func TestOpen(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	tu.RequireNotNil(t, store)
}

func TestCreate(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	err := store.Create(post.Post{Id: "sqlite-test"})
	tu.RequireNil(t, err)

	posts, _ := store.FindAll()
	tu.RequireEqual(t, len(posts), 1)

	post, _ := store.FindById("sqlite-test")
	tu.RequireNotNil(t, post)
}
