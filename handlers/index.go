package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func IndexGet(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join(".", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
