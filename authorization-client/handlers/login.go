package handlers

import (
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	html := `<a href="/start-oauth">Fazer login com Authorization Server</a>`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
