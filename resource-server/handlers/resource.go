package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"resource-server/store"
	"resource-server/token"
)

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Printf("Authorization header missing")
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		log.Printf("Invalid Authorization header format: %s", authHeader)
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	authToken := authHeader[len(bearerPrefix):]
	claims, err := token.Validate(authToken)
	if err != nil {
		log.Printf("Error validating token: %v", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		log.Printf("Claim 'sub' not found: %v", claims)
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	events := store.GetEventsByUserId(userId)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(events); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
