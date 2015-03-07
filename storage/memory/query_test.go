package memory

import (
	"testing"
	"time"

	"../../post"
	storage ".."
	tu "../../util/testing"
)

var examplePosts = []post.Post{
	post.Post{"1", "first post", "something important!", time.Now()},
	post.Post{"2", "second post", "a realization.", time.Now()},
}

func TestFindOnlyId(t *testing.T) {
	store := FromPosts(examplePosts)

	// find a post
	postFindById, _ := store.FindById("1")
	q, _ := storage.Query().Find("id", "1").Build()
	postFind, _ := store.Find(*q)

	tu.RequireEqual(t, len(postFind), 1)
	tu.ExpectEqual(t, *postFindById, postFind[0])
}

func TestNotFindOnlyId(t *testing.T) {
	store := FromPosts(examplePosts)

	// don't find a post
	postFindById, _ := store.FindById("3")
	q, _ := storage.Query().Find("id", "3").Build()
	postFind, _ := store.Find(*q)

	tu.ExpectEqual(t, len(postFind), 0)
	tu.ExpectEqual(t, postFindById, (*post.Post)(nil))
}

func TestFindAll(t *testing.T) {
	store := FromPosts(examplePosts)

	postsFindAll, _ := store.FindAll()
	q, _ := storage.Query().Build()
	postsFind, _ := store.Find(*q)

	tu.ExpectEqual(t, len(postsFindAll), len(postsFind))
}

func TestFindByTitle(t *testing.T) {
	store := FromPosts(examplePosts)

	q, _ := storage.Query().Find("title", "first post").Build()
	postsFind, err := store.Find(*q)

	tu.RequireNil(t, err)
	tu.RequireEqual(t, len(postsFind), 1)
	tu.ExpectEqual(t, postsFind[0].Title, "first post")


	q, _ = storage.Query().Find("title", "second post").Build()
	postsFind, err = store.Find(*q)

	tu.RequireNil(t, err)
	tu.RequireEqual(t, len(postsFind), 1)
	tu.ExpectEqual(t, postsFind[0].Title, "second post")
	tu.ExpectEqual(t, postsFind[0].Content, "a realization.")
}

func TestFindByIdWithQuery(t *testing.T) {
	store := FromPosts(examplePosts)

	q, _ := storage.Query().Find("id", "1").Build()
	postsFind, err := store.runFind(*q)

	tu.RequireNil(t, err)
	tu.RequireEqual(t, len(postsFind), 1)
	tu.ExpectEqual(t, postsFind[0].Id, "1")
	tu.ExpectEqual(t, postsFind[0].Title, "first post")

}
