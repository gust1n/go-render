package main

import (
	"log"
	"net/http"

	"github.com/gust1n/go-render/render"
)

var rnd *render.Renderer

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := rnd.ExecuteTemplate(w, "index.html", map[string]interface{}{"Title": "Home"}); err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if err := rnd.ExecuteTemplate(w, "profile.html", map[string]interface{}{"Title": "Profile"}); err != nil {
			log.Println(err)
		}
	})
	http.HandleFunc("/map", func(w http.ResponseWriter, r *http.Request) {
		if err := rnd.ExecuteTemplate(w, "map.html", map[string]interface{}{"Title": "Map"}); err != nil {
			log.Println(err)
		}
	})
	log.Println("web server listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	rnd = render.New()

	// Render all the pages
	// The extend recursion is automatically made so /pages is enough
	// You could also pass just templates but then you have to remember to change to pages/index.html and so on above
	rnd.LoadTemplates("templates/pages")
}
