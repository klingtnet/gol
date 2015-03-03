package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
)

func main() {
	var homePageTemplate = template.Must(template.New("homepage").Parse(homePageTemplateStr))
	var content = blackfriday.MarkdownCommon([]byte(`# gol

## subheading

- I
- am
- a
- list

[source](https://github.com/KLINGTdotNET/gol)`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m := make(map[string]interface{})
		m["title"] = "gol"
		m["message"] = "Hello, World!"
		m["content"] = template.HTML(content)
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

		<p>{{ .message }}</p>

		{{ .content }}	
	</body>
</html>`
