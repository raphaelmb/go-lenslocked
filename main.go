package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/raphaelmb/go-lenslocked/controllers"
	"github.com/raphaelmb/go-lenslocked/templates"
	"github.com/raphaelmb/go-lenslocked/views"
)

func main() {
	r := chi.NewRouter()

	tmpl := views.Must(views.ParseFS(templates.FS, "home.tmpl.html", "tailwind.tmpl.html"))
	r.Get("/", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "contact.tmpl.html", "tailwind.tmpl.html"))
	r.Get("/contact", controllers.StaticHandler(tmpl))

	tmpl = views.Must(views.ParseFS(templates.FS, "faq.tmpl.html", "tailwind.tmpl.html"))
	r.Get("/faq", controllers.FAQ(tmpl))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
