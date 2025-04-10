package handlers

import "net/http"

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("response_type") != "code" {
		//TODO handle token response_type
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	redirect := r.URL.Query().Get("redirect_uri")

	w.Header().Set("Location", redirect+"?code=example-code")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
