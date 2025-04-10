package handlers

import (
	"authorization-client/config"
	"fmt"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf(`
						<!DOCTYPE html>
						<html lang="en">
						<head>
						<meta charset="UTF-8">
						<meta name="viewport" content="width=device-width, initial-scale=1.0">
						<title>Login</title>
						</head>
						<body>
						<a href="%s?response_type=code&client_id=abcde&redirect_uri=%s&scope=read_data&state=random_state">Fazer login externo</a>
						</body>
						</html>`, config.AuthorizationServerUrl("authorize", false), config.Url("callback"))

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}
