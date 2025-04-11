package handlers

import (
	"authorization-client/config"
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(
		`<a href="%s?response_type=%s&client_id=%s&redirect_uri=%s&scope=%s&state=%s">Fazer login com Authorization Server</a>`,
		config.AuthorizationServerUrl("authorize", false),
		"code",
		config.ClientID(),
		config.Url("callback"),
		"read write",
		"random_state",
	)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
