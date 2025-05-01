package handlers

import (
	"authorization-client/config"
	"authorization-client/store"
	"errors"
	"log"
	"net/http"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var html strings.Builder

	html.WriteString("<h1>Client</h1>")

	_, err := getToken(r)
	if err != nil {
		html.WriteString(`<p>Welcome, to access the Resource Server, click <a href="/auth_server/initiate">here</a> to authorize.</p>`)
	} else {
		html.WriteString(`<p>Welcome, you are logged in</p>`)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html.String()))
}

func getToken(r *http.Request) (*store.Tokens, error) {
	session, err := r.Cookie(config.SESSION_COOKIE_NAME)
	if err != nil || session == nil || session.Value == "" {
		log.Printf("error: session cookie not found, error: %v", err)
		return nil, errors.New("session cookie not found")
	}

	tokens, err := store.GetToken(session.Value)
	if err != nil {
		log.Printf("error: failed to get token from store, error: %v", err)
		return nil, errors.New("failed to get token from store")
	}

	return tokens, nil
}
