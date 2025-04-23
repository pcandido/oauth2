package handlers

import (
	"authorization-client/config"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"
)

func LoginInitiateHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState(16)

	http.SetCookie(w, stateCookie(state))
	http.Redirect(w, r, redirectUrl(state), http.StatusTemporaryRedirect)
}

func redirectUrl(state string) string {
	query := url.Values{}
	query.Add("response_type", "code")
	query.Add("client_id", config.ClientID())
	query.Add("redirect_uri", config.Url("login/callback"))
	query.Add("scope", "read write")
	query.Add("state", state)

	return config.AuthorizationServerUrl("authorize", true, &query)
}

func stateCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:   "oauth_state",
		Value:  state,
		MaxAge: int((5 * time.Minute).Seconds()),
	}
}

func generateRandomState(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // Handle error appropriately in a real application
	}
	return hex.EncodeToString(bytes)
}
