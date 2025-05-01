package handlers

import (
	"authorization-client/config"
	"authorization-client/utils"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func LoginInitiateHandler(w http.ResponseWriter, r *http.Request) {
	state := utils.GenerateRandomString(16)

	http.SetCookie(w, stateCookie(state))
	http.Redirect(w, r, redirectUrl(state), http.StatusTemporaryRedirect)
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
	query.Add("redirect_uri", fmt.Sprintf("%s/auth_server/callback", config.BaseUrl()))
	query.Add("scope", "read:events")
	query.Add("state", state)

	url := &url.URL{
		Scheme:   "http",
		Host:     config.AuthServerHostForFrontend(),
		Path:     "authorize",
		RawQuery: query.Encode(),
	}

	return url.String()
}
