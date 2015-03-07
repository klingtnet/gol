package memory

import (
	"errors"

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
	} else {
		return nil, errors.New("queries not supported")
	}
}

