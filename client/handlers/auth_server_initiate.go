package handlers

import (
	"authorization-client/config"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func LoginInitiateHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState(16)

	http.SetCookie(w, stateCookie(state))
	http.Redirect(w, r, redirectUrl(state), http.StatusTemporaryRedirect)
}

func generateRandomState(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func stateCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:     config.AUTH_SERVER_STATE_COOKIE_NAME,
		Value:    state,
		MaxAge:   int((5 * time.Minute).Seconds()),
		HttpOnly: true,
		// Use secure flag in production
	}
}

func redirectUrl(state string) string {
	query := url.Values{}
	query.Add("response_type", "code")
	query.Add("client_id", config.CLIENT_ID)
	query.Add("redirect_uri", fmt.Sprintf("%s/login/callback", config.BaseUrl()))
	query.Add("scope", "read write")
	query.Add("state", state)

	url := &url.URL{
		Scheme:   "http",
		Host:     config.AuthServerHostForFrontend(),
		Path:     "authorize",
		RawQuery: query.Encode(),
	}

	return url.String()
}
