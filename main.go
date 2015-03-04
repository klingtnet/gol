package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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

	sanitizePolicy := bluemonday.UGCPolicy()
	sanitizePolicy.AllowElements("iframe", "audio", "video")
	sanitizePolicy.AllowAttrs("width", "height", "src").OnElements("iframe", "audio", "video", "img")
	templateUtils := template.FuncMap{
		"markdown": func(content string) template.HTML {
			htmlContent := blackfriday.MarkdownCommon([]byte(content))
			htmlContent = sanitizePolicy.SanitizeBytes(htmlContent)
			return template.HTML(htmlContent)
		},
		"formatTime": func(t time.Time) template.HTML {
			// thanks, http://fuckinggodateformat.com/ (every language/template thingy should have this)
			isoDate := t.Format(time.RFC3339)
			readableDate := t.Format("January 2, 2006 (15:04)")
			return template.HTML(fmt.Sprintf("<time datetime=\"%s\">%s</time>", isoDate, readableDate))
		},
		"assetUrl": func(path string) string {
			return fmt.Sprintf("%s/%s", assetBase, path)
		},
	}
	homePageTemplate := template.New("homepage").Funcs(templateUtils)
	homePageTemplate = template.Must(homePageTemplate.Parse(homePageTemplateStr))

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
		homePageTemplate.Execute(w, m)
	})

	createPostTemplate := template.New("create").Funcs(templateUtils)
	createPostTemplate = template.Must(createPostTemplate.Parse(createPostTemplateStr))

	router.HandleFunc("/posts/new", func(w http.ResponseWriter, r *http.Request) {
		createPostTemplate.Execute(w, map[string]string{"title": "Write a new post!"})
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
			createPostTemplate.Execute(w, m)
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

var homePageTemplateStr = `<!DOCTYPE html>
<html lang=en>
	<head>
		<title>{{ .title }}</title>

		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/css/materialize.min.css">

		<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/styles/tomorrow.min.css">

		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"/>
		<link rel="stylesheet" href="{{ "main.css" | assetUrl }}" />
	</head>

	<body>
		<div class="container">
			<div id="edit-button" class="fixed-action-btn">
				<a href="/posts/new" class="btn-floating btn-large waves-effect waves-light blue tooltipped" data-tooltip="Write a new post"><i class="mdi-content-add"></i></a>
			</div>

			{{ range $post := .posts }}
			<article id="post-{{ $post.Id }}" class="post">
				<div class="post-actions">
					<a href="/posts/{{ $post.Id }}/edit" class="btn-floating waves-effect waves-light blue tooltipped" data-tooltip="Edit post"><i class="mdi-editor-mode-edit"></i></a>
					<a href="/posts/{{ $post.Id }}" data-method="DELETE" class="btn-floating waves-effect waves-light red tooltipped" data-tooltip="Delete post"><i class="mdi-action-delete"></i></a>
				</div>
				<h1><a href="/posts/{{ $post.Id }}">{{ $post.Title }}</a></h1>
				<h5>Posted on <i>{{ $post.Created | formatTime }}</i></h5>

				<div class="post-content flow-text">
					{{ $post.Content | markdown }}
				</div>
			</article>
			<hr />
			{{ end }}
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdn.rawgit.com/heyLu/materialize.css/master/dist/js/materialize.min.js"></script>

		<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/highlight.min.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>

		<script src="{{ "main.js" | assetUrl }}"></script>
	</body>
</html>`

var createPostTemplateStr = `<!DOCTYPE html>
<html lang=en>
	<head>
		<title>{{ .title }}</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/css/materialize.min.css">
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"/>
		<link rel="stylesheet" href="{{ "main.css" | assetUrl }}" />
	</head>

	<body>
		<div class="container">
			<h1>{{ .title }}</h1>

			<form method="POST" action="/posts{{ if .post }}/{{ .post.Id }}{{ end }}">
				<div class="input-field">
					<input class="markdown-input" name="title" type="text" value="{{ .post.Title }}"></input>
					<label for="title">Titlemania</label>
				</div>
				<div class="input-field">
					<textarea class="materialize-textarea markdown-input" name="content" rows="80" cols="100">{{ .post.Content }}</textarea>
					<label for="content">Your thoughts.</label>
				</div>


				<button class="btn waves-effect waves-light" type="submit" name="action">
					Submit
				</button>
			</form>
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdn.rawgit.com/heyLu/materialize.css/master/dist/js/materialize.min.js"></script>

		<script src="{{ "main.js" | assetUrl }}"></script>
	</body>
</html>`
