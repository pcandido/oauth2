package handlers

import (
	"authorization-server/store"
	"authorization-server/token"
	"encoding/json"
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
		http.Error(w, "invalid_request", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	grantType := r.FormValue("grant_type")
	code := r.FormValue("code")
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	clientSecret := r.FormValue("client_secret")

	if grantType != "authorization_code" || code == "" || redirectURI == "" || clientID == "" {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	authCode, err := store.GetCode(code, clientID, redirectURI)
	if err != nil {
		http.Error(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	client, err := store.GetClient(clientID)
	if err != nil {
		http.Error(w, "invalid_client", http.StatusBadRequest)
		return
	}
	if client.ClientSecret != clientSecret {
		http.Error(w, "invalid_client", http.StatusBadRequest)
		return
	}

	user, err := store.GetUserById(authCode.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, err := token.Generate(map[string]any{
		"sub":   user.ID,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"aud":   clientID,
		"scope": authCode.Scope,
	}, 15*time.Minute)
	if err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}

	refreshToken, err := token.Generate(map[string]any{
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"aud": clientID,
	}, 7*time.Hour*24)
	if err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int((15 * time.Minute).Seconds()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
