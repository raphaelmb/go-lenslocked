package controllers

import (
	"net/http"

	"github.com/raphaelmb/go-lenslocked/views"
)

func StaticHandler(tmpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}
