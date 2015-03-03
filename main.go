package main

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
	postsJson, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, postsJson, 0644)
	if err != nil {
		return err
	}

	return nil
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

func main() {
	posts, err := readPosts("posts.json")
	if err != nil {
		log.Fatal(err)
	}

	templateUtils := template.FuncMap{
		"markdown": func(content string) template.HTML {
			htmlContent := blackfriday.MarkdownCommon([]byte(content))
			return template.HTML(htmlContent)
		},
		"formatTime": func(t time.Time) template.HTML {
			// thanks, http://fuckinggodateformat.com/ (every language/template thingy should have this)
			isoDate := t.Format(time.RFC3339)
			readableDate := t.Format("January 2, 2006 (15:04)")
			return template.HTML(fmt.Sprintf("<time datetime=\"%s\">%s</time>", isoDate, readableDate))
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

	createPostTemplate := template.Must(template.New("create").Parse(createPostTemplateStr))

	router.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		createPostTemplate.Execute(w, nil)
	})

	router.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // POST creates a new post
			now := time.Now()
			post := Post{
				Id:      fmt.Sprintf("%x", md5.Sum(toByteSlice(now))),
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
				Created: now,
			}
			posts, _ = readPosts("posts.json")
			posts = append(posts, post)
			writePosts("posts.json", posts)
			json.NewEncoder(w).Encode(post)
		} else { // TODO: GET list all posts
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("not implemented"))
		}
	})

	// http.HandleFunc("/posts", ...) // GET = display all posts
	// http:HandleFunc("/posts/:id", ...) // GET/POST = get/edit an existing post

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
		<style>
			#edit-button {
				top: 23px;
			}

			.post .post-actions {
				visibility: hidden;
				float: right;
			}

			.post:hover .post-actions {
				visibility: visible;
			}

			.post .post-content {
				font-size: 1.414rem;
				line-height: 1.8rem;
			}
		</style>
	</head>

	<body>
		<div class="container">
			<div id="edit-button" class="fixed-action-btn">
				<a href="/create" class="btn-floating btn-large waves-effect waves-light blue"><i class="mdi-content-add"></i></a>
			</div>

			{{ range $post := .posts }}
			<article class="post">
				<div class="post-actions">
					<a href="/edit" class="btn-floating waves-effect waves-light blue"><i class="mdi-editor-mode-edit"></i></a>
					<a href="/edit" class="btn-floating waves-effect waves-light red"><i class="mdi-action-delete"></i></a>
				</div>
				<h1>{{ $post.Title }}</h1>
				<h5>Posted on <i>{{ $post.Created | formatTime }}</i></h5>

				<div class="post-content flow-text">
					{{ $post.Content | markdown }}
				</div>
			</article>
			<hr />
			{{ end }}
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/js/materialize.min.js"></script>

		<script src="//cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/highlight.min.js"></script>
		<script>hljs.initHighlightingOnLoad();</script>
	</body>
</html>`

var createPostTemplateStr = `<!DOCTYPE html>
<html lang=en>
	<head>
		<title>{{ .title }}</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/css/materialize.min.css">
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"/>
	</head>

	<body>
		<div class="container">
			<h1>Write a new post!</h1>

			<form method="POST" action="/posts">
				<div class="input-field">
					<input name="title" type="text"></input>
					<label for="title">Titlemania</label>
				</div>
				<div class="input-field">
					<textarea class="materialize-textarea" name="content" rows="50" cols="120"></textarea>
					<label for="content">Your thoughts.</label>
				</div>


				<button class="btn waves-effect waves-light" type="submit" name="action">
					Submit
				</button>
			</form>
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/js/materialize.min.js"></script>
	</body>
</html>`
