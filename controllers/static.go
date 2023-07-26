package controllers

import (
	"html/template"
	"net/http"

	"github.com/raphaelmb/go-lenslocked/views"
)

func StaticHandler(tmpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func FAQ(tmpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes, there is a free version.",
		},
		{
			Question: "What are you support hours?",
			Answer:   "We are available 24/7.",
		},
		{
			Question: "How do I contact support",
			Answer:   `You can reach us via email - <a href=\"mailto:email@example.com\">email@example.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, questions)
	}
}
