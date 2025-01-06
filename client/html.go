package client

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

var (
	//go:embed templates/**
	templatesFS embed.FS
	tmpl        *template.Template
)

func init() {
	var err error
	tmpl, err = template.ParseFS(templatesFS, "templates/**")
	if err != nil {
		panic(fmt.Errorf("error parsing templates: %w", err))
	}
}

func renderAuthPage(w http.ResponseWriter, d authData) {
	if err := tmpl.ExecuteTemplate(w, "authentication_response.html", d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}
