package handlers

import (
	"authorization-client/config"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

func StartOauthHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomString(16)

	redirectUrl := fmt.Sprintf(
		`%s?response_type=%s&client_id=%s&redirect_uri=%s&scope=%s&state=%s`,
		config.AuthorizationServerUrl("authorize", false),
		"code",
		config.ClientID(),
		config.Url("callback"),
		"read write",
		state,
	)

	fmt.Println(state)

	http.SetCookie(w, &http.Cookie{
		Name:  "oauth_state",
		Value: state,
		// use Secure and HttpOnly flags for production
	})

	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err) // Handle error appropriately in a real application
	}
	return hex.EncodeToString(bytes)
}
