package sqlite

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"testing"
	"time"

	storage ".."
	"../../post"
	tu "../../util/testing"
)

func makePost(id, title, content string) post.Post {
	return post.Post{
		Id:      id,
		Created: time.Now(),
		Title:   title,
		Content: content,
	}
}

func comparePosts(t *testing.T, a, b *post.Post) {
	tu.RequireEqual(t, a.Id, b.Id)
	tu.RequireEqual(t, a.Created.Unix(), b.Created.Unix())
	tu.RequireEqual(t, a.Title, b.Title)
	tu.RequireEqual(t, a.Content, b.Content)
}

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

func TestFind(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	post := makePost("0815", "sqlite-test-find", "Testing store.Find()")
	store.Create(post)

	//TODO: needs the query interface
}

func TestFindById(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	post := makePost("0815", "sqlite-test-find-by-id", "Testing store.FindById()")
	store.Create(post)

	foundPost, err := store.FindById(post.Id)
	tu.RequireNil(t, err)
	comparePosts(t, foundPost, &post)
}

func TestFindAll(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	posts := make([]post.Post, 3)
	for i, _ := range posts {
		post := makePost(fmt.Sprintf("%d", i*i), fmt.Sprintf("post-#%d", i), fmt.Sprintf("This is post number %d", i))
		posts[i] = post
		store.Create(post)
	}

	foundPosts, err := store.FindAll()
	tu.RequireNil(t, err)
	tu.RequireEqual(t, len(foundPosts), len(posts))

	for _, foundPost := range foundPosts {
		for _, post := range posts {
			if foundPost.Id == post.Id {
				comparePosts(t, &foundPost, &post)
			}
		}
	}
}

func TestCreate(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	post := makePost("0815", "sqlite-test-create", "Testing sqlite post creation")
	err := store.Create(post)
	tu.RequireNil(t, err)

	posts, _ := store.FindAll()
	tu.RequireEqual(t, len(posts), 1)

	foundPost, _ := store.FindById("0815")

	tu.RequireNotNil(t, foundPost)
	comparePosts(t, foundPost, &post)
}

func TestUpdate(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	post := makePost("0815", "sqlite-test-update", "Testing sqlite post updayte")
	updatedPost := makePost("0815", "sqlite-test-update", "Testing sqlite post update")

	store.Create(post)
	store.Update(updatedPost)
	foundPost, _ := store.FindById("0815")
	tu.RequireNotEqual(t, foundPost.Content, post.Content)
	comparePosts(t, foundPost, &updatedPost)
}

func TestDelete(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	store.Create(makePost("0815", "sqlite-test-delete", "A one-time post."))
	err := store.Delete("0815")
	tu.RequireNil(t, err)
	posts, _ := store.FindAll()
	tu.RequireEqual(t, len(posts), 0)
}

func TestClose(t *testing.T) {
	store, tearDown := tSetup(t)
	defer tearDown()

	store.Close()
	err := store.Create(makePost("0815", "sqlite-test-delete", "Post creation on closed storage must fail!"))
	tu.RequireNotNil(t, err)
}
