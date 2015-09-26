package sqlite

import (
	//	"errors"
	"fmt"
	"time"

	"../../post"
	"../query"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func (s *Store) Find(q query.Query) ([]post.Post, error) {
	posts := make([]post.Post, 0)
	var query string
	var err error = nil
	var rows *sql.Rows
	if q.Find != nil {
		query = fmt.Sprintf("SELECT * FROM posts WHERE %s = ? ORDER BY ? DESC", q.Find.Name)
		stmt, err := s.db.Prepare(query)
		if err != nil {
			return nil, err
		}

		rows, err = stmt.Query(q.Find.Value, q.SortBy)
	} else {
		query = "SELECT * FROM posts ORDER BY ? DESC"
		stmt, err := s.db.Prepare(query)
		if err != nil {
			return nil, err
		}

		rows, err = stmt.Query(q.SortBy)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var title, content, id string
		var created time.Time
		err = rows.Scan(&id, &created, &title, &content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post.Post{id, title, content, created})
	}

	return posts, nil
}
