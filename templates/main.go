package templates

import (
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"time"
)

func Templates(assetBase string) *template.Template {
	sanitizePolicy := bluemonday.UGCPolicy()
	sanitizePolicy.AllowElements("iframe", "audio", "video")
	sanitizePolicy.AllowAttrs("width", "height", "src").OnElements("iframe", "audio", "video", "img")

	templateFuncs := template.FuncMap{
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

	templateTree := template.New("").Funcs(templateFuncs)

	// shared templates (must be loaded first to be availlable in the
	// page templates)
	template.Must(templateTree.ParseGlob("templates/shared/*.tmpl"))

	// page templates
	template.Must(templateTree.ParseGlob("templates/*.tmpl"))

	return templateTree
}
