package sqlite

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/url"
	"time"

	storage ".."
	"../../post"
)

type Backend struct{}

type Store struct {
	path string
	db   *sql.DB
}

func init() {
	storage.Register("sqlite", Backend{})
}

func setup(db *sql.DB) error {
	// initialize sqlite database if not present under path
	creatTableStmt := "CREATE TABLE IF NOT EXISTS posts (id TEXT NOT NULL PRIMARY KEY, created DATETIME, title TEXT, content TEXT)"
	// db.Exec does not return results
	_, err := db.Exec(creatTableStmt)
	if err != nil {
		return err
	}

	createIndexStmt := "CREATE UNIQUE INDEX IF NOT EXISTS idIdx ON posts (id)"
	_, err = db.Exec(createIndexStmt)
	return err
}

func (m Backend) Open(u *url.URL) (storage.Store, error) {
	path := u.Host + u.Path
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = setup(db)
	if err != nil {
		log.Fatal(err)
	}

	// return store
	store := storage.Store(&Store{
		db: db,
	})
	return store, nil
}

// Store interface methods
func (s *Store) FindById(id string) (*post.Post, error) {
	stmt, err := s.db.Prepare("SELECT * FROM posts WHERE ID = ?")
	if err != nil {
		return nil, err
	}
	// never returns nil
	row := stmt.QueryRow(id)

	var title, content string
	var created time.Time
	err = row.Scan(&id, &created, &title, &content)
	post := &post.Post{id, title, content, created}

	switch {
	case err == sql.ErrNoRows:
		return nil, errors.New("could not find post with id!")
	case err != nil:
		return nil, err
	}

	return post, nil
}

func (s *Store) FindAll() ([]post.Post, error) {
	rows, err := s.db.Query("SELECT id, created, title, content FROM posts ORDER BY created DESC")
	if err != nil {
		return nil, err
	}

	posts := make([]post.Post, 0)
	defer rows.Close()
	for rows.Next() {
		var id, title, content string
		var created time.Time
		err = rows.Scan(&id, &created, &title, &content)
		if err != nil {
			log.Print(err)
		}
		posts = append(posts, post.Post{id, title, content, created})
	}
	return posts, rows.Err()
}

func (s *Store) Create(post post.Post) error {
	return s.execQuery("INSERT INTO posts(id, created, title, content) values(?, ?, ?, ?)", post.Id, post.Created, post.Title, post.Content)
}

func (s *Store) Update(updatedPost post.Post) error {
	_, err := s.FindById(updatedPost.Id)
	if err != nil {
		return err
	}

	return s.execQuery("UPDATE posts SET id=?, created=?, title=?, content=? WHERE id=?", updatedPost.Id, updatedPost.Created, updatedPost.Title, updatedPost.Content, updatedPost.Id)
}

func (s *Store) Delete(id string) error {
	return s.execQuery("DELETE FROM posts WHERE id = ?", id)
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Sync() error {
	// TODO
	return nil
}

func (s *Store) execQuery(query string, args ...interface{}) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("could not begin transaction", err)
		return err
	}

	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Println("could not prepare query", err)
		return err
	}

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Println("could not execute statement", err)
		return err
	}

	return tx.Commit()
}
