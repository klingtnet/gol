package post

import (
	"sort"
	"time"
)

type Post struct {
	Id      string    `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
}

type ByDate []Post

func (p ByDate) Len() int           { return len(p) }
func (p ByDate) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByDate) Less(i, j int) bool { return p[i].Created.Unix() < p[j].Created.Unix() }

func Sort(sortable sort.Interface) sort.Interface {
	sort.Sort(sortable)
	return sortable
}

func Reverse(sortable sort.Interface) sort.Interface {
	sort.Sort(sort.Reverse(sortable))
	return sortable
}


