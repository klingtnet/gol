package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"./templates"
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

func readPosts(filename string) ([]Post, error) {
	var posts []Post
	postsJson, err := ioutil.ReadFile("posts.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(postsJson, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func writePosts(filename string, posts []Post) error {
	postsJson, err := json.MarshalIndent(posts, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, postsJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func findPost(posts []Post, id string) *Post {
	for i, post := range posts {
		if post.Id == id {
			return &posts[i]
		}
	}

	return nil
}

func deletePost(posts []Post, id string) ([]Post, error) {
	newPosts := make([]Post, 0, len(posts))
	foundPost := false

	for _, post := range posts {
		if post.Id != id {
			newPosts = append(newPosts, post)
		} else {
			foundPost = true
		}
	}

	if !foundPost {
		return posts, errors.New("post not found")
	}

	return newPosts, nil
}

func toByteSlice(data interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		log.Println("Error: binary.Write failed:", err)
		return []byte{}
	}
	return buf.Bytes()
}

func notImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

var Environment = "development"
var Version = "master"
var assetBase = "/assets"

func init() {
	if Environment == "production" {
		assetBase = fmt.Sprintf("https://cdn.rawgit.com/KLINGTdotNET/gol/%s/assets", Version)
	}

	fmt.Printf("gol - v%s (%s)\n", Version, Environment)
}

func main() {
	posts, err := readPosts("posts.json")
	if err != nil {
		log.Fatal(err)
	}

	templates := templates.Templates(assetBase)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		posts, err = readPosts("posts.json")
		if err != nil {
			log.Println("Warning: Could not read posts.json:", err)
		}
		m := make(map[string]interface{})
		m["title"] = "gol"
		m["posts"] = Reverse(ByDate(posts))
		templates.ExecuteTemplate(w, "posts", m)
	})

	router.HandleFunc("/posts/new", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "post_form", map[string]string{"title": "Write a new post!"})
	})

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // POST creates a new post
			now := time.Now()
			post := Post{
				Id:      fmt.Sprintf("%x", md5.Sum(toByteSlice(now.UnixNano()))),
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
				Created: now,
			}
			posts, _ = readPosts("posts.json")
			posts = append(posts, post)
			writePosts("posts.json", posts)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else { // TODO: GET list all posts
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		post := findPost(posts, id)
		if post == nil {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		if r.Method == "GET" {
			if post != nil {
				json.NewEncoder(w).Encode(post)
			}
		} else if r.Method == "HEAD" {
			// already handle by post == nil above
		} else if r.Method == "POST" {
			var newPost Post
			if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				newPost.Title = r.FormValue("title")
				newPost.Content = r.FormValue("content")

				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else { // assume it's JSON
				err := json.NewDecoder(r.Body).Decode(&newPost)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
				}
				w.WriteHeader(http.StatusAccepted)
			}

			if newPost.Title != "" {
				post.Title = newPost.Title
			}
			if newPost.Content != "" {
				post.Content = newPost.Content
			}
			writePosts("posts.json", posts)
			json.NewEncoder(w).Encode(post)
		} else if r.Method == "DELETE" {
			posts, err = deletePost(posts, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			writePosts("posts.json", posts)
		} else {
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		post := findPost(posts, id)
		if post != nil {
			m := make(map[string]interface{})
			m["title"] = "Edit post"
			m["post"] = post
			templates.ExecuteTemplate(w, "post_form", m)
		} else {
			http.NotFound(w, r)
		}
	})

	// http.HandleFunc("/posts", ...) // GET = display all posts

	if Environment == "development" {
		// in development, serve local assets
		router.PathPrefix("/assets").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	}

	http.Handle("/", router)

	fmt.Println("Listening on http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
