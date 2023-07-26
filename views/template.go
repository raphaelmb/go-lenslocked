package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tmpl, err := template.ParseFS(fs, patterns...)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tmpl}, nil
}

func Parse(filepath string) (Template, error) {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tmpl}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
