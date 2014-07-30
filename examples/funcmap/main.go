package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gust1n/go-render/render"
)

var templates map[string]*template.Template

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templates["index.html"].Execute(w, map[string]interface{}{"Name": "Joakim"}); err != nil {
			log.Println(err)
		}
	})
	log.Println("web server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	var tmplErr error
	templates, tmplErr = render.LoadWithFuncMap("templates", template.FuncMap{
		"greet": func(name string) string {
			return fmt.Sprintf("Hello %s", name)
		},
	})
	if tmplErr != nil {
		panic(tmplErr)
	}
}
