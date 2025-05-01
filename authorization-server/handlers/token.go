package handlers

import (
	"authorization-server/store"
	"authorization-server/token"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("invalid request method: %s", r.Method)
		http.Error(w, "invalid_request", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	if grantType != "authorization_code" {
		log.Printf("invalid grant type: %s", grantType)
		http.Error(w, "invalid_grant_type", http.StatusBadRequest)
		return
	}

	if code == "" {
		log.Printf("missing code")
		http.Error(w, "invalid_code", http.StatusBadRequest)
		return
	}

	authorization, err := store.PopAuthorization(code)
	if err != nil {
		log.Printf("error getting authorization \"%s\": %v", code, err)
		http.Error(w, "invalid_grant", http.StatusUnauthorized)
		return
	}

	if authorization.ClientId != clientID {
		log.Printf("client ID mismatch: expected %s, got %s", authorization.ClientId, clientID)
		http.Error(w, "invalid_grant", http.StatusUnauthorized)
		return
	}

	if authorization.RedirectUri != redirectURI {
		log.Printf("redirect URI mismatch: expected %s, got %s", authorization.RedirectUri, redirectURI)
		http.Error(w, "invalid_grant", http.StatusUnauthorized)
		return
	}

	client, err := store.GetClient(clientID)
	if err != nil {
		log.Printf("error getting client \"%s\": %v", clientID, err)
		http.Error(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	if client.ClientSecret != clientSecret {
		log.Printf("client secret mismatch")
		http.Error(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	user, err := store.GetUserById(authorization.UserId)
	if err != nil {
		log.Printf("error getting user \"%s\": %v", authorization.UserId, err)
		http.Error(w, "invalid_user", http.StatusUnauthorized)
		return
	}

	accessTokenExpiresIn := 15 * time.Minute
	accessToken, err := token.Generate(map[string]any{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"aud":   clientID,
		"scope": authorization.Scope,
	}, accessTokenExpiresIn)

	if err != nil {
		log.Printf("error generating access token: %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := token.Generate(map[string]any{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"aud": clientID,
	}, 7*time.Hour*24)

	if err != nil {
		log.Printf("error generating refresh token: %v", err)
		http.Error(w, "internal_server_error", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(accessTokenExpiresIn.Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
