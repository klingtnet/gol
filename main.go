package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"./post"
	"./storage"
	_ "./storage/memory"
	"./templates"
)

func readPosts(filename string) ([]post.Post, error) {
	var posts []post.Post
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

func writePosts(filename string, posts []post.Post) error {
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

func findPost(posts []post.Post, id string) *post.Post {
	for i, post := range posts {
		if post.Id == id {
			return &posts[i]
		}
	}

	return nil
}

func deletePost(posts []post.Post, id string) ([]post.Post, error) {
	newPosts := make([]post.Post, 0, len(posts))
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

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func notImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

var Environment = getEnv("ENVIRONMENT", "development")
var Version = "master"
var assetBase = "/assets"
var ssl = flag.String("ssl", "", "enable ssl (give server.crt,server.key as value)")

func init() {
	if Environment == "production" {
		assetBase = fmt.Sprintf("https://cdn.rawgit.com/KLINGTdotNET/gol/%s/assets", Version)
	}

	fmt.Printf("gol - v%s (%s)\n", Version, Environment)
}

func main() {
	flag.Parse()

	store, _ := storage.Open("memory://")
	fmt.Println(store)

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
		m["posts"] = post.Reverse(post.ByDate(posts))
		templates.ExecuteTemplate(w, "posts", m)
	})

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		posts, _ := readPosts("posts.json")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}).Methods("GET").Headers("Content-Type", "application/json")

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // POST creates a new post
			now := time.Now()
			post := post.Post{
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

	router.HandleFunc("/posts/new", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "post_form", map[string]string{"title": "Write a new post!"})
	})

	router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		p := findPost(posts, id)
		if p == nil {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		if r.Method == "GET" {
			if p != nil {
				json.NewEncoder(w).Encode(p)
			}
		} else if r.Method == "HEAD" {
			// already handle by p == nil above
		} else if r.Method == "POST" {
			var newPost post.Post
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
				p.Title = newPost.Title
			}
			if newPost.Content != "" {
				p.Content = newPost.Content
			}
			writePosts("posts.json", posts)
			json.NewEncoder(w).Encode(p)
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

	port := getEnv("PORT", "5000")
	if *ssl == "" {
		fmt.Printf("Listening on http://0.0.0.0:%s\n", port)
		log.Fatal(http.ListenAndServe(":" + port, nil))
	} else {
		certAndKey := strings.Split(*ssl, ",")
		if len(certAndKey) != 2 {
			fmt.Println("Error: -ssl needs server.crt,server.key as arguments")
			os.Exit(1)
		}
		fmt.Printf("Listening on https://0.0.0.0:%s\n", port)
		log.Fatal(http.ListenAndServeTLS(":" + port, certAndKey[0], certAndKey[1], nil))
	}
}
