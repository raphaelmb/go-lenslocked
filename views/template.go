package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/csrf"
	"github.com/raphaelmb/go-lenslocked/context"
	"github.com/raphaelmb/go-lenslocked/models"
)

type public interface {
	Public() string
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tmpl := template.New(path.Base(patterns[0]))
	tmpl = tmpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("csrfField not implemented")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
		"errors": func() []string {
			return nil
		},
	})

	tmpl, err := tmpl.ParseFS(fs, patterns...)
	if err != nil {
		log.Printf("error parsing template: %v", err)
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{htmlTpl: tmpl}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	tmpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("error cloning template: %v", err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
	errMsgs := errMessages(errs...)
	tmpl = tmpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			return errMsgs
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

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
