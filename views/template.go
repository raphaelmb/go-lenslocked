package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/raphaelmb/go-lenslocked/context"
	"github.com/raphaelmb/go-lenslocked/models"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tmpl := template.New(patterns[0])
	tmpl = tmpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
	})

	tmpl, err := tmpl.ParseFS(fs, patterns...)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{htmlTpl: tmpl}, nil
}

// func Parse(filepath string) (Template, error) {
// 	tmpl, err := template.ParseFiles(filepath)
// 	if err != nil {
// 		log.Printf("error parsing template: %v", err)
// 		return Template{}, fmt.Errorf("parsing template: %w", err)
// 	}

// 	return Template{htmlTpl: tmpl}, nil
// }

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	tmpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("error cloning template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
	tmpl = tmpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}
