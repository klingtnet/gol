// not a storage, but the query interface
package query

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Field struct {
	Name  string
	Value interface{}
}

type Query struct {
	Find       *Field
	Start      int // negative == not specified
	Count      int
	Matches    []Field
	RangeStart *time.Time
	RangeEnd   *time.Time
	SortBy     string
	Reverse    bool
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

func New() Builder {
	return &DefaultBuilder{Default}
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
	if err := valueIn("match", field, []string{"id", "title", "content"}); err != nil {
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
	b.query.SortBy = field
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

func (q Invalid) Find(field string, value interface{}) Builder  { return q }
func (q Invalid) Start(pos uint) Builder                        { return q }
func (q Invalid) Count(pos uint) Builder                        { return q }
func (q Invalid) Match(field string, value interface{}) Builder { return q }
func (q Invalid) Range(start, end time.Time) Builder            { return q }
func (q Invalid) SortBy(field string) Builder                   { return q }
func (q Invalid) Reverse() Builder                              { return q }

func (q Invalid) Build() (*Query, error) {
	return nil, q.Err
}

// build from query params
//
// Start(10).Count(30) == ?start=10&count=30
// Range(now, twoDaysAgo) == ?range=now:twoDaysAgo
// SortBy("title") == ?sort=title
// Reverse() == ?reverse
// Matches("title", "cool") == ?match=title:cool
// Matches("title", "cool").Matches("content", "wow") == ?match=title:cool&match=content:cool
func FromParams(params url.Values) (*Query, error) {
	b := New()

	for key, vals := range params {
		//fmt.Println(key, vals)

		v := vals[0]
		switch key {
		case "id":
			b = b.Find("id", vals[len(vals)-1])
		case "title":
			b = b.Find("title", vals[len(vals)-1])
		case "start":
			start, err := parsePos("start", v)
			if err != nil {
				return nil, err
			}
			b = b.Start(start)
		case "count":
			count, err := parsePos("count", v)
			if err != nil {
				return nil, err
			}
			b = b.Count(count)
		case "sort":
			b = b.SortBy(v)
		case "reverse":
			if v == "" || v == "true" {
				b = b.Reverse()
			}
		case "match":
			for _, m := range vals {
				matchPair := strings.Split(m, ":")
				if len(matchPair) != 2 {
					return nil, errors.New(fmt.Sprintf("match must be of the format field:match, but was '%s'", m))
				}
				b = b.Match(matchPair[0], matchPair[1])
			}
		case "range":
			rangePair := strings.Split(v, ",")
			if len(rangePair) != 2 {
				return nil, errors.New(fmt.Sprintf("range must be of the format `start,end`, but was '%s'", v))
			}
			start, err := time.Parse(time.RFC3339, rangePair[0])
			if err != nil {
				return nil, errors.New(fmt.Sprint("invalid range start: ", err))
			}
			end, err := time.Parse(time.RFC3339, rangePair[1])
			if err != nil {
				return nil, errors.New(fmt.Sprint("invalid range and: ", err))
			}
			b = b.Range(start, end)
		}
	}

	return b.Build()
}

func parsePos(name, s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, errors.New(fmt.Sprint("invalid %s value:", name, err))
	}
	return uint(i), nil
}

// TODO: func ToParams(q Query) url.Values
