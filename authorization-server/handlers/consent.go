package handlers

import (
	"authorization-server/config"
	"authorization-server/store"
	"authorization-server/token"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func ConsentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		consentForm(w, r)
		return
	}

	if r.Method == http.MethodPost {
		handleConsent(w, r)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func consentForm(w http.ResponseWriter, r *http.Request) {
	authorizeParams := AuthorizeParamsFromQuery(r)

	client, err := store.GetClient(authorizeParams.ClientID)
	if err != nil {
		http.Error(w, "Client not found", http.StatusBadRequest)
		return
	}

	htmlResponse := fmt.Sprintf(`
            <form method="POST" action="/consent">
                <input type="hidden" name="client_id" value="%s">
                <input type="hidden" name="redirect_uri" value="%s">
                <input type="hidden" name="response_type" value="%s">
                <input type="hidden" name="scope" value="%s">
                <input type="hidden" name="state" value="%s">

                <p>The application <strong>%s</strong> is requesting access to your data with the following scopes:</p>
                <ul>
                    %s
                </ul>
                <button type="submit" name="action" value="approve">Approve</button>
                <button type="submit" name="action" value="deny">Deny</button>
                <button type="submit" name="action" value="logout">Logout</button>
            </form>
        `,
		html.EscapeString(authorizeParams.ClientID),
		html.EscapeString(authorizeParams.RedirectURI),
		html.EscapeString(authorizeParams.ResponseType),
		html.EscapeString(authorizeParams.Scope),
		html.EscapeString(authorizeParams.State),
		html.EscapeString(client.ClientName),
		formatScopes(authorizeParams.Scope),
	)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlResponse))
}

func handleConsent(w http.ResponseWriter, r *http.Request) {
	authorizeParams := AuthorizeParamsFromForm(r)
	action := r.FormValue("action")

	if action == "logout" {
		http.SetCookie(w, &http.Cookie{
			Name:     config.ACCESS_TOKEN_COOKIE_NAME,
			Value:    "",
			HttpOnly: true,
			MaxAge:   -1,
		})

		loginUrl := url.URL{
			Path:     "/login",
			RawQuery: authorizeParams.toQuery().Encode(),
		}

		http.Redirect(w, r, loginUrl.String(), http.StatusFound)
		return
	}

	if action == "deny" {
		authorizeParams.Error = "access_denied"
		url := url.URL{
			Path:     "/authorize",
			RawQuery: authorizeParams.toQuery().Encode(),
		}
		http.Redirect(w, r, url.String(), http.StatusFound)
		return
	}

	accessToken, err := r.Cookie(config.ACCESS_TOKEN_COOKIE_NAME)
	if err != nil {
		log.Printf("error getting access token cookie: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, err := token.Validate(accessToken.Value)
	if err != nil {
		log.Printf("error validating access token: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userId, ok := claims["user"].(string)
	if !ok {
		log.Printf("error getting user ID from claims: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	consentToken, err := token.Generate(map[string]any{
		"user":         userId,
		"client_id":    authorizeParams.ClientID,
		"redirect_uri": authorizeParams.RedirectURI,
		"scope":        authorizeParams.Scope,
	}, 5*time.Minute)

	if err != nil {
		log.Printf("error generating consent token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	authorizeParams.ConsentToken = consentToken
	url := url.URL{
		Path:     "/authorize",
		RawQuery: authorizeParams.toQuery().Encode(),
	}

	http.Redirect(w, r, url.String(), http.StatusFound)
}

func formatScopes(scopes string) string {
	scopeList := ""
	for _, scope := range splitScopes(scopes) {
		scopeList += fmt.Sprintf("<li>%s</li>", html.EscapeString(scope))
	}
	return scopeList
}

func splitScopes(scopes string) []string {
	return strings.Fields(scopes)
}
