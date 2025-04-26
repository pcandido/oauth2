package handlers

import (
	"authorization-server/config"
	"authorization-server/store"
	"authorization-server/token"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginForm(w, r)
		return
	}

	if r.Method == http.MethodPost {
		login(w, r)
		return
	}
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	authorizeParams := AuthorizeParamsFromQuery(r)

	htmlResponse := fmt.Sprintf(`
			<form method="POST" action="/login">
				<input type="hidden" name="client_id" value="%s">
				<input type="hidden" name="redirect_uri" value="%s">
				<input type="hidden" name="response_type" value="%s">
				<input type="hidden" name="scope" value="%s">
				<input type="hidden" name="state" value="%s">

				<label for="email">Email:</label>
				<input type="text" id="email" name="email" required />
				<label for="password">Password:</label>
				<input type="password" id="password" name="password" required />
				<button type="submit">Login</button>
			</form>
		`,
		html.EscapeString(authorizeParams.ClientID),
		html.EscapeString(authorizeParams.RedirectURI),
		html.EscapeString(authorizeParams.ResponseType),
		html.EscapeString(authorizeParams.Scope),
		html.EscapeString(authorizeParams.State),
	)

	if authorizeParams.Error != "" {
		htmlResponse += fmt.Sprintf("<p style='color:red;'>%s</p>", html.EscapeString(authorizeParams.Error))
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlResponse))
}

func login(w http.ResponseWriter, r *http.Request) {
	authorizeParams := AuthorizeParamsFromForm(r)

	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := store.GetUserByEmail(email)
	if err != nil || !user.ValidatePassword(password) {
		authorizeParams.Error = "invalid_credentials"

		url := url.URL{
			Path:     "/login",
			RawQuery: authorizeParams.toQuery().Encode(),
		}

		http.Redirect(w, r, url.String(), http.StatusFound)
		return
	}

	token, err := token.Generate(map[string]any{"user": user.ID}, 15*time.Minute)
	if err != nil {
		authorizeParams.Error = "server_error"
		url := url.URL{
			Path:     "/authorize",
			RawQuery: authorizeParams.toQuery().Encode(),
		}
		http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
		return
	}

	cookie := &http.Cookie{
		Name:  config.ACCESS_TOKEN_COOKIE_NAME,
		Value: token,
		// Use Secure and HttpOnly flags for security in production
	}
	http.SetCookie(w, cookie)

	url := url.URL{
		Path:     "/authorize",
		RawQuery: authorizeParams.toQuery().Encode(),
	}

	http.Redirect(w, r, url.String(), http.StatusTemporaryRedirect)
}
