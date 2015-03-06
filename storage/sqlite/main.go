package sqlite

import (
	"database/sql"
	"errors"
	//	"fmt"
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
	creatTableStmt := "CREATE TABLE IF NOT EXISTS posts (id TEXT NOT NULL PRIMARY KEY, created DATETIME, title TEXT, content TEXT)"
	// db.Exec does not return results
	_, err := db.Exec(creatTableStmt)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (m Backend) Open(u *url.URL) (storage.Store, error) {
	path := u.Host + u.Path
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// initialize sqlite database if not present under path
	setup(db)
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

	// ugghhh
	var id1, title, content string
	var created time.Time
	err = row.Scan(&id1, &created, &title, &content)
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
	rows, err := s.db.Query("SELECT id, created, title, content FROM posts")
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
	query := "INSERT INTO posts(id, created, title, content) values(?, ?, ?, ?)"

	tx, err := s.db.Begin()
	if err != nil {
		log.Print("could not begin transaction!")
		return err
	}
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Print("could not prepare statement!", err)
		return err
	}
	_, err = stmt.Exec(post.Id, post.Created, post.Title, post.Content)
	if err != nil {
		log.Printf("could not execute statement!")
		return err
	}
	return tx.Commit()
}

func (s *Store) Update(updatedPost post.Post) error {
	_, err := s.FindById(updatedPost.Id)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := s.db.Prepare("UPDATE posts SET id=?, created=?, title=?, content=? WHERE id=?")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = stmt.Exec(updatedPost.Id, updatedPost.Created, updatedPost.Title, updatedPost.Content, updatedPost.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) Delete(id string) error {
	query := "DELETE FROM posts WHERE id = ?"

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Sync() error {
	// TODO
	return nil
}
