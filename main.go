package main

import (
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
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
	}
	homePageTemplate := template.New("homepage").Funcs(templateUtils)
	homePageTemplate = template.Must(homePageTemplate.Parse(homePageTemplateStr))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := make(map[string]interface{})
		m["title"] = "gol"
		m["posts"] = posts
		homePageTemplate.Execute(w, m)
	})

	createPostTemplate := template.Must(template.New("create").Parse(createPostTemplateStr))

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		createPostTemplate.Execute(w, nil)
	})

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" { // POST creates a new post
			post := Post{
				Title:   r.FormValue("title"),
				Content: r.FormValue("content"),
			}
			postJson, _ := json.Marshal(post)
			w.Write(postJson)
		} else { // TODO: GET list all posts
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte("not implemented"))
		}
	})

	// http.HandleFunc("/posts", ...) // GET = display all posts
	// http:HandleFunc("/posts/:id", ...) // GET/POST = get/edit an existing post

	fmt.Println("Listening on http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

var homePageTemplateStr = `<!DOCTYPE html>
<html lang=en>
	<head>
		<title>{{ .title }}</title>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/css/materialize.min.css">
		<meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no"/>
	</head>

	<body>
		<div class="container">
			<p>All posts...</p>

			{{ range $post := .posts }}
				<h1>{{ $post.Title }}</h1>

				{{ $post.Content | markdown }}
			{{ end }}
		</div>

		<script type="text/javascript" src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.95.3/js/materialize.min.js"></script>
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
