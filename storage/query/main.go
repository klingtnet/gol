// not a storage, but the query interface
package query

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Field struct {
	Name string
	Value interface{}
}

type Query struct{
	Find *Field
	Start int // negative == not specified
	Count int
	Matches []Field
	RangeStart *time.Time
	RangeEnd *time.Time
	SortBy string
	Reverse bool
}

// default is to get all posts, sorted by created date
var Default = Query{nil, -1, -1, nil, nil, nil, "created", false}

func IsDefault(q Query) bool {
	return q.Find == nil && q.Start == -1 && q.Count == -1 && q.Matches == nil &&
		q.RangeStart == nil && q.RangeEnd == nil && q.SortBy == "created" && !q.Reverse
}

type Builder interface {
	Find(field string, value interface{}) Builder // exact match
	Start(pos uint) Builder
	Count(pos uint) Builder
	Match(field string, value interface{}) Builder // partial match
	Range(start, end time.Time) Builder
	SortBy(field string) Builder
	Reverse() Builder
	Build() (*Query, error)
}

type DefaultBuilder struct {
	query Query
}

func (b *DefaultBuilder) Find(field string, value interface{}) Builder {
	if err := valueIn("by", field, []string{"id", "title", "created"}); err != nil {
		return Invalid{err}
	}
	b.query.Find = &Field{field, value}
	return b
}

func (b *DefaultBuilder) Start(pos uint) Builder {
	b.query.Start = int(pos)
	return b
}

func (b *DefaultBuilder) Count(count uint) Builder {
	b.query.Count = int(count)
	return b
}

func (b *DefaultBuilder) Match(field string, value interface{}) Builder {
	if err := valueIn("match", field, []string{"id", "title"}); err != nil {
		return Invalid{err}
	}
	b.query.Matches = append(b.query.Matches, Field{field, value})
	return b
}

func (b *DefaultBuilder) Range(start, end time.Time) Builder {
	if start.Unix() == end.Unix() {
		return Invalid{errors.New("empty range")}
	}
	b.query.RangeStart = &start
	b.query.RangeEnd = &end
	return b
}

func (b *DefaultBuilder) SortBy(field string) Builder {
	if err := valueIn("sort", field, []string{"title", "created"}); err != nil {
		return Invalid{err}
	}
	return b
}

func (b *DefaultBuilder) Reverse() Builder {
	b.query.Reverse = !b.query.Reverse
	return b
}

func valueIn(name string, value string, values []string) error {
	for _, v := range values {
		if v == value {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("%s must be one of %#v but was %#v", name, values, value))
}

func (b *DefaultBuilder) Build() (*Query, error) {
	return &b.query, nil
}

// supports further chaining, returns the error; used by backends
// (returned when an invalid query has been detected)
//
// backends may alternatively chose to ignore parts of the query
type Invalid struct {
	Err error
}

func (q Invalid) Find(field string, value interface{}) Builder { return q }
func (q Invalid) Start(pos uint) Builder { return q }
func (q Invalid) Count(pos uint) Builder { return q }
func (q Invalid) Match(field string, value interface{}) Builder { return q }
func (q Invalid) Range(start, end time.Time) Builder { return q }
func (q Invalid) SortBy(field string) Builder { return q }
func (q Invalid) Reverse() Builder { return q }

func (q Invalid) Build() (*Query, error) {
	return nil, q.Err
}

// build from query params
func FromParams(params url.Values) (*Query, error) {
	for key, val := range params {
		fmt.Println(key, val)
	}

	return nil, nil
}
