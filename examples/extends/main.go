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
		if err := templates["index.html"].Execute(w, map[string]interface{}{"Title": "Home"}); err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if err := templates["profile.html"].Execute(w, map[string]interface{}{"Title": "Profile"}); err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		if err := templates["map.html"].Execute(w, map[string]interface{}{"Title": "Map"}); err != nil {
			log.Println(err)
		}
	})
	log.Println("web server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	// Pre-render/Load all the pages
	// The extend recursion is automatically made so /pages is enough
	// You could also pass just templates but then you have to remember to change to pages/index.html and so on above
	var tmplErr error
	if templates, tmplErr = render.Load("templates/pages"); tmplErr != nil {
		panic(tmplErr)
	}
}
