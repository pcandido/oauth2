package handlers

import (
	"net/http"
	"resource-server/token"
)

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	authToken := authHeader[len(bearerPrefix):]
	_, err := token.Validate(authToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[
		{
			"title": "Meeting with client",
			"start": "2025-05-01T10:00:00Z",
			"end": "2025-05-01T12:00:00Z"
		},
		{
			"title": "Product review",
			"start": "2025-05-01T13:00:00Z",
			"end": "2025-05-01T15:00:00Z"
		},
		{
			"title": "Retrospective",
			"start": "2025-05-01T16:00:00Z",
			"end": "2025-05-01T16:00:00Z"
		}
	]`))
}
