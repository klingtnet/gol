package memory

import (
	"errors"
	"fmt"

	"../../post"
	"../query"
)

func (s *Store) Find(q query.Query) ([]post.Post, error) {
	if q.Find != nil && q.Find.Name == "id" {
		id, ok := q.Find.Value.(string)
		if !ok {
			return nil, errors.New("id must be a string")
		}

		p, err := s.FindById(id)
		if err != nil {
			return nil, err
		}
		return []post.Post{*p}, nil
	} else if query.IsDefault(q) {
		return s.FindAll()
	} else if q.Find != nil {
		return s.runFind(q)
	} else if q.Start != -1 || q.Count != -1 {
		return s.runQuery(q)
	} else {
		return nil, errors.New("query not supported")
	}
}


func (s *Store) runFind(q query.Query) ([]post.Post, error) {
	for _, p := range s.posts {
		found := false

		switch q.Find.Name {
		case "id":
			found = p.Id == q.Find.Value
		case "title":
			found = p.Title == q.Find.Value
		default:
			return nil, errors.New(fmt.Sprint("unsupported field:", q.Find.Name))
		}

		if found {
			return []post.Post{p}, nil
		}
	}

	return []post.Post{}, nil
}

func (s *Store) runQuery(q query.Query) ([]post.Post, error) {
	start := 0
	if q.Start != -1 {
		start = q.Start
	}
	count := len(s.posts)
	if q.Count != -1 {
		count = q.Count
	}

	// TODO: adjust capacity based on query type (probably not)
	posts := make([]post.Post, 0, 10)
	if count == 0 || start >= len(s.posts) {
		return posts, nil
	}

	n := 0
	for _, post := range s.posts {
		if n >= start {
			posts = append(posts, post)
		}

		n += 1

		if n >= start + count {
			return posts, nil
		}
	}

	return posts, nil
}
