package handlers

import (
	"authorization-server/config"
	"authorization-server/store"
	"authorization-server/token"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	authorizeParams := AuthorizeParamsFromQuery(r)
	if err := validateClient(authorizeParams); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		// TODO return errors to client
		return
	}

	if authorizeParams.Error != "" {
		http.Error(w, authorizeParams.Error, http.StatusBadRequest)
		return
	}

	if !logedIn(r) {
		loginUrl := url.URL{
			Path:     "/login",
			RawQuery: authorizeParams.toQuery().Encode(),
		}

		http.Redirect(w, r, loginUrl.String(), http.StatusTemporaryRedirect)
		return
	}

	callbackQuery := url.Values{}
	callbackQuery.Set("code", "example-code")
	//TODO generate a real code
	callbackQuery.Set("state", authorizeParams.State)

	callbackUrl, err := url.Parse(authorizeParams.RedirectURI)
	if err != nil {
		http.Error(w, "invalid redirect_uri", http.StatusBadRequest)
		return
	}

	callbackUrl.RawQuery = callbackQuery.Encode()
	http.Redirect(w, r, callbackUrl.String(), http.StatusTemporaryRedirect)
}

func logedIn(r *http.Request) bool {
	cookie, err := r.Cookie(config.ACCESS_TOKEN_COOKIE_NAME)
	if err != nil || cookie.Value == "" {
		return false
	}

	_, err = token.Validate(cookie.Value)

	return err == nil
}

func validateClient(authorizeParams *AuthorizeParams) error {
	client, err := store.GetClient(authorizeParams.ClientID)
	if err != nil {
		return fmt.Errorf("unauthorized_client")
	}

	if !slices.Contains(client.RedirectURIs, authorizeParams.RedirectURI) {
		return fmt.Errorf("invalid_redirect_uri")
	}

	if !slices.Contains(client.ResponseTypes, authorizeParams.ResponseType) {
		return fmt.Errorf("unsupported_response_type")
	}

	scopes := strings.Split(authorizeParams.Scope, " ")
	if len(scopes) == 0 {
		return fmt.Errorf("invalid_scope")
	}

	for _, scope := range scopes {
		if !slices.Contains(client.Scopes, scope) {
			return fmt.Errorf("invalid_scope")
		}
	}

	return nil
}
