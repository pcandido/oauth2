package handlers

import (
	"authorization-server/config"
	"authorization-server/store"
	"authorization-server/token"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	authorizeParams := AuthorizeParamsFromQuery(r)

	if err := validateClient(authorizeParams); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	callbackUrl, err := url.Parse(authorizeParams.RedirectURI)
	if err != nil {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}

	callbackQuery := url.Values{}
	callbackQuery.Set("state", authorizeParams.State)

	if authorizeParams.Error != "" {
		callbackQuery.Set("error", authorizeParams.Error)
		callbackUrl.RawQuery = callbackQuery.Encode()

		http.Redirect(w, r, callbackUrl.String(), http.StatusTemporaryRedirect)
		return
	}

	userId, loggedIn := getUser(r)
	if !loggedIn {
		loginUrl := url.URL{
			Path:     "/login",
			RawQuery: authorizeParams.toQuery().Encode(),
		}

		http.Redirect(w, r, loginUrl.String(), http.StatusTemporaryRedirect)
		return
	}

	if authorizeParams.ConsentToken == "" {
		consentUrl := url.URL{
			Path:     "/consent",
			RawQuery: authorizeParams.toQuery().Encode(),
		}

		http.Redirect(w, r, consentUrl.String(), http.StatusTemporaryRedirect)
		return
	}

	consentClaims, err := token.Validate(authorizeParams.ConsentToken)
	if err != nil {
		log.Printf("error validating consent token: %v", err)
		http.Error(w, "invalid consent token", http.StatusBadRequest)
		return
	}

	consentUserId, ok := consentClaims["user"].(string)
	if !ok || consentUserId != userId {
		log.Printf("user ID mismatch: expected %s, got %s", userId, consentUserId)
		http.Error(w, "invalid consent token", http.StatusBadRequest)
		return
	}
	consentClientId, ok := consentClaims["client_id"].(string)
	if !ok || consentClientId != authorizeParams.ClientID {
		log.Printf("client ID mismatch: expected %s, got %s", authorizeParams.ClientID, consentClientId)
		http.Error(w, "invalid consent token", http.StatusBadRequest)
		return
	}
	consentRedirectURI, ok := consentClaims["redirect_uri"].(string)
	if !ok || consentRedirectURI != authorizeParams.RedirectURI {
		log.Printf("redirect URI mismatch: expected %s, got %s", authorizeParams.RedirectURI, consentRedirectURI)
		http.Error(w, "invalid consent token", http.StatusBadRequest)
		return
	}
	consentScope, ok := consentClaims["scope"].(string)
	if !ok || consentScope != authorizeParams.Scope {
		log.Printf("scope mismatch: expected %s, got %s", authorizeParams.Scope, consentScope)
		http.Error(w, "invalid consent token", http.StatusBadRequest)
		return
	}

	code := store.PushAuthorization(userId, authorizeParams.ClientID, authorizeParams.RedirectURI, authorizeParams.Scope)
	callbackQuery.Set("code", code)

	callbackUrl.RawQuery = callbackQuery.Encode()
	http.Redirect(w, r, callbackUrl.String(), http.StatusTemporaryRedirect)
}

func getUser(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(config.ACCESS_TOKEN_COOKIE_NAME)
	if err != nil || cookie.Value == "" {
		return "", false
	}

	claims, err := token.Validate(cookie.Value)
	if err != nil {
		return "", false
	}

	userId, ok := claims["user"].(string)
	if !ok {
		return "", false
	}

	return userId, true
}

func validateClient(authorizeParams *AuthorizeParams) error {
	client, err := store.GetClient(authorizeParams.ClientID)
	if err != nil {
		log.Printf("error getting client \"%s\": %v", authorizeParams.ClientID, err)
		return fmt.Errorf("unauthorized_client")
	}

	if !slices.Contains(client.RedirectURIs, authorizeParams.RedirectURI) {
		log.Printf("redirect URI mismatch: expected %s, got %s", client.RedirectURIs, authorizeParams.RedirectURI)
		return fmt.Errorf("invalid_redirect_uri")
	}

	if !slices.Contains(client.ResponseTypes, authorizeParams.ResponseType) {
		log.Printf("unsupported response type: expected %v, got %s", client.ResponseTypes, authorizeParams.ResponseType)
		return fmt.Errorf("unsupported_response_type")
	}

	scopes := strings.Split(authorizeParams.Scope, " ")
	if len(scopes) == 0 {
		log.Printf("invalid scope: %s", authorizeParams.Scope)
		return fmt.Errorf("invalid_scope")
	}

	for _, scope := range scopes {
		if !slices.Contains(client.Scopes, scope) {
			log.Printf("invalid scope. client options: %v, scope: %s", client.Scopes, scope)
			return fmt.Errorf("invalid_scope")
		}
	}

	return nil
}
