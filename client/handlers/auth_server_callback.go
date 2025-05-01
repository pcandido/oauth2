package handlers

import (
	"authorization-client/config"
	"authorization-client/store"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	callbackState := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	// Check if the state is valid
	cookieState, err := r.Cookie(config.AUTH_SERVER_STATE_COOKIE_NAME)
	if err != nil {
		log.Printf("error: state cookie not found, error: %v", err)
		http.Error(w, "state cookie not found", http.StatusBadRequest)
		return
	}

	if cookieState.Value != callbackState {
		log.Printf("error: invalid state, expected: %s, got: %s", cookieState.Value, callbackState)
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   config.AUTH_SERVER_STATE_COOKIE_NAME,
		Value:  "",
		MaxAge: -1,
	})

	// Exchange the authorization code for an access token
	tokenUrl := url.URL{
		Scheme: "http",
		Host:   config.AuthServerHostForBackend(),
		Path:   "token",
	}

	tokenBody := url.Values{}
	tokenBody.Set("grant_type", "authorization_code")
	tokenBody.Set("code", code)
	tokenBody.Set("redirect_uri", fmt.Sprintf("%s/auth_server/callback", config.BaseUrl()))
	tokenBody.Set("client_id", config.CLIENT_ID)
	tokenBody.Set("client_secret", config.CLIENT_SECRET)

	res, err := http.Post(tokenUrl.String(), "application/x-www-form-urlencoded", strings.NewReader(tokenBody.Encode()))
	if err != nil {
		log.Printf("error: failed to exchange authorization code for token, error: %v", err)
		http.Error(w, "error getting token", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("error: token endpoint returned non-200 status, status: %d", res.StatusCode)
		http.Error(w, "error getting token", res.StatusCode)
		return
	}

	var tokens store.Tokens
	if err := json.NewDecoder(res.Body).Decode(&tokens); err != nil {
		log.Printf("error: failed to parse token response, error: %v", err)
		http.Error(w, "error parsing token response", http.StatusInternalServerError)
		return
	}

	session := store.StoreToken(tokens)

	http.SetCookie(w, &http.Cookie{
		Name:     config.SESSION_COOKIE_NAME,
		Value:    session,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		// Use secure flag in production
	})

	// Redirect to the main page
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
