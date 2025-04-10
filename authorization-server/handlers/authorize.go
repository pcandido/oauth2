package handlers

import (
	"authorization-server/store"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")
	scopes := strings.Split(r.URL.Query().Get("scope"), " ")

	client, err := store.GetClient(clientID)
	if err != nil {
		http.Error(w, "Client not found", http.StatusBadRequest) //TODO transformar as validações em redirect com erro
		return
	}

	if !slices.Contains(client.RedirectURIs, redirectURI) {
		http.Error(w, "Invalid redirect_uri", http.StatusBadRequest)
		return
	}

	if !slices.Contains(client.ResponseTypes, responseType) {
		http.Error(w, "Invalid response_type", http.StatusBadRequest)
		return
	}

	for _, scope := range scopes {
		if !slices.Contains(client.Scopes, scope) {
			fmt.Printf("Invalid scope: %s\n", scopes)
			http.Error(w, "Invalid scope", http.StatusBadRequest)
			return
		}
	}

	//TODO validar scope

	w.Header().Set("Location", fmt.Sprintf("%s?code=example-code&state=%s", redirectURI, r.URL.Query().Get("state")))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
