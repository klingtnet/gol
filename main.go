package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
)

type Post struct {
	Title   string
	Content string
}

func main() {
	post := Post{
		Title: "My First Post!",
		Content: `# gol

## subheading

- I
- am
- a
- list

[source](https://github.com/KLINGTdotNET/gol)`}

	var homePageTemplate = template.Must(template.New("homepage").Parse(homePageTemplateStr))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := make(map[string]interface{})
		m["title"] = post.Title
		htmlContent := blackfriday.MarkdownCommon([]byte(post.Content))
		m["content"] = template.HTML(htmlContent)
		homePageTemplate.Execute(w, m)
	})

	fmt.Println("Listening on http://0.0.0.0:5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

var homePageTemplateStr = `<!DOCTYPE html>
<html lang=en>
	<head>
		<title>{{ .title }}</title>
	</head>

	<body>
		<h1>{{ .title }}</h1>

		{{ .content }}	
	</body>
</html>`
