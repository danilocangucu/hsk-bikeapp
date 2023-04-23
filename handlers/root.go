package handlers

import (
	"net/http"
	"regexp"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	if match, _ := regexp.MatchString(`^/.+$`, r.URL.Path); match {
		http.NotFound(w, r)
		return
	} else {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	}
}
