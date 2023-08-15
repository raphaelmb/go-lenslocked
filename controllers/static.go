package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tmpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, r, nil)
	}
}

func FAQ(tmpl Template) http.HandlerFunc {
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
		tmpl.Execute(w, r, questions)
	}
}
