package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gust1n/go-render/render"
)

var templates map[string]*template.Template

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templates["pages/index.html"].Execute(w, map[string]interface{}{"Title": "Home"}); err != nil {
			log.Println(err)
		}
	})
	log.Println("web server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	var tmplErr error
	if templates, tmplErr = render.Load("templates"); tmplErr != nil {
		panic(tmplErr)
	}
}
