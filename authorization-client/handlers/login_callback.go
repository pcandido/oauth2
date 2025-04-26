package handlers

import (
	"authorization-client/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func LoginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if err := validateState(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	deleteStateCookie(w)

	tokens, err := getTokensFromCode(r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeTokenCookies(w, tokens)

	w.Header().Set("Location", config.Url(""))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func validateState(r *http.Request) error {
	cookie, err := r.Cookie(config.OAUTH_STATE_COOKIE_NAME)
	if err != nil {
		return fmt.Errorf("state cookie not found: %v", err)
	}

	if cookie.Value != r.URL.Query().Get("state") {
		return fmt.Errorf("invalid state")
	}

	return nil
}

func deleteStateCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   config.OAUTH_STATE_COOKIE_NAME,
		Value:  "",
		MaxAge: -1,
	})
}

func getTokensFromCode(code string) (*Tokens, error) {
	body := fmt.Sprintf(
		"grant_type=authorization_code&code=%s&redirect_uri=%s&client_id=%s&client_secret=%s",
		code,
		config.Url("login/callback"),
		config.ClientID(),
		config.ClientSecret(),
	)

	tokenUrl := config.AuthorizationServerUrl("token", false, nil)

	res, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error getting token: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting token: %d", res.StatusCode)
	}

	var tokens Tokens
	if err := json.NewDecoder(res.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("error parsing token response: %s", err)
	}

	return &tokens, nil
}

func writeTokenCookies(w http.ResponseWriter, tokens *Tokens) {
	http.SetCookie(w, &http.Cookie{
		Name:  "access_token",
		Value: tokens.AccessToken,
		Path:  "/",
		// For production, set Secure and HttpOnly
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "refresh_token",
		Value: tokens.RefreshToken,
		Path:  "/",
		// For production, set Secure and HttpOnly
	})
}
