package main

import (
	"bytes"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pmylund/go-cache"
	"github.com/zenazn/goji"
	gojiweb "github.com/zenazn/goji/web"
)

var tmplt *template.Template
var storage *cache.Cache

func main() {
	rand.Seed(time.Now().UnixNano())
	loadTemplates()
	storage = cache.New(time.Hour, 5*time.Minute)

	static := gojiweb.New()
	static.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/static/", static)

	goji.Get("/", Index)
	goji.Get("/team", Team)
	goji.Get("/team/:team_id", Team)
	goji.Get("/solo", Solo)

	goji.Serve()
}

func loadTemplates() {
	var templates []string

	fn := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() != true && strings.HasSuffix(f.Name(), ".html") {
			templates = append(templates, path)
		}
		return nil
	}

	err := filepath.Walk("./views/", fn)

	if err != nil {
		panic(err.Error())
	}

	tmplt = template.Must(template.ParseFiles(templates...))
}

func parseTemplate(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}
