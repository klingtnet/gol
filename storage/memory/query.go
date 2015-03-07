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
	} else {
		return nil, errors.New("queries not supported")
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
