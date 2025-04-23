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
	state := generateRandomString(16)

	query := url.Values{}
	query.Add("response_type", "code")
	query.Add("client_id", config.ClientID())
	query.Add("redirect_uri", config.Url("login/callback"))
	query.Add("scope", "read write")
	query.Add("state", state)

	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  state,
		MaxAge: int((5 * time.Minute).Seconds()),
		// use Secure and HttpOnly flags for production
	})

	http.Redirect(w, r, config.AuthorizationServerUrl("authorize", true, &query), http.StatusTemporaryRedirect)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // Handle error appropriately in a real application
	}
	return hex.EncodeToString(bytes)
}
