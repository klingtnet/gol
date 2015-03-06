package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ogier/pflag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"./post"
	"./storage"
	_ "./storage/memory"
	_ "./storage/json"
	"./templates"
)

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

func renderPosts(templates *template.Template, w http.ResponseWriter, posts []post.Post) {
	m := make(map[string]interface{})
	m["title"] = "gol"
	m["posts"] = post.Reverse(post.ByDate(posts))
	templates.ExecuteTemplate(w, "posts", m)
}

func notImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

var Environment = getEnv("ENVIRONMENT", "development")
var Version = "master"
var assetBase = "/assets"
var ssl = pflag.String("ssl", "", "enable ssl (give server.crt,server.key as value)")
var storageUrl = pflag.String("storage", "json://posts.json", "the storage to connect to")

func init() {
	if Environment == "production" {
		assetBase = fmt.Sprintf("https://cdn.rawgit.com/KLINGTdotNET/gol/%s/assets", Version)
	}

	fmt.Printf("gol - %s (%s)\n", Version, Environment)
}

func main() {
	pflag.Parse()

	store, err := storage.Open(*storageUrl)
	if err != nil {
		log.Fatal(err)
	}

	templates := templates.Templates(assetBase)

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		posts, err := store.FindAll()
		if err != nil {
			log.Println("Warning: Could not read posts.json:", err)
		}
		renderPosts(templates, w, posts)
	})

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		posts, _ := store.FindAll()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}).Methods("GET").Headers("Content-Type", "application/json")

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			posts, _ := store.FindAll()
			renderPosts(templates, w, posts)
		} else if r.Method == "POST" { // POST creates a new post
			now := time.Now()
			post := post.Post{
				Id:      fmt.Sprintf("%x", md5.Sum(toByteSlice(now.UnixNano()))),
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
				Created: now,
			}

			err := store.Create(post)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/new", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "post_form", map[string]string{"title": "Write a new post!"})
	})

	router.HandleFunc("/posts/preview", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // use post to receive content in body
			var post post.Post
			err := json.NewDecoder(r.Body).Decode(&post)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				templates.ExecuteTemplate(w, "post", post)
			}
		} else {
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		p, _ := store.FindById(id)
		if p == nil {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}

		if r.Method == "GET" {
			if r.Header.Get("Content-Type") == "application/json" {
				json.NewEncoder(w).Encode(p)
			} else {
				m := make(map[string]interface{})
				m["title"] = p.Title
				m["posts"] = []post.Post{*p}
				templates.ExecuteTemplate(w, "posts", m)
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
			store.Update(*p)
			json.NewEncoder(w).Encode(p)
		} else if r.Method == "DELETE" {
			err := store.Delete(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
			}
		} else {
			notImplemented(w)
		}
	})

	router.HandleFunc("/posts/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		post, _ := store.FindById(id)
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
