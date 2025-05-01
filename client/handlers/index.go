package handlers

import (
	"authorization-client/config"
	"net/http"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	var html strings.Builder

	html.WriteString("<h1>Client</h1>")

	cookie, err := r.Cookie(config.AUTH_SERVER_ACCESS_TOKEN_COOKIE_NAME)
	if err != nil || cookie.Value == "" {
		html.WriteString(`<p>Bem vindo, para acessar o Resource Server, clique <a href="/auth_server/initiate">aqui</a> para autorizar.</p>`)
	} else {
		html.WriteString(`<p>Bem vindo, vc está dentro</p>`)
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html.String()))
}
