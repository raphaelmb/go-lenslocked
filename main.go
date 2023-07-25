package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/raphaelmb/go-lenslocked/controllers"
	"github.com/raphaelmb/go-lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	tmpl, err := views.Parse(filepath.Join("templates", "home.tmpl.html"))
	if err != nil {
		panic(err)
	}
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl, err = views.Parse(filepath.Join("templates", "contact.tmpl.html"))
	if err != nil {
		panic(err)
	}
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl, err = views.Parse(filepath.Join("templates", "faq.tmpl.html"))
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.StaticHandler(tmpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
