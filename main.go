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
	Title   string `json:title`
	Content string `json:content`
}

func main() {
	var posts []Post
	postsJson, err := ioutil.ReadFile("posts.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(postsJson, &posts)

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
	</body>
</html>`
