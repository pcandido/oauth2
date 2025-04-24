package main

import (
	"authorization-server/handlers"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/authorize", handlers.AuthorizeHandler)
	http.HandleFunc("/token", handlers.TokenHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	// http.HandleFunc("/revoke", revokeHandler)
	// http.HandleFunc("/introspect", introspectHandler)
	// http.HandleFunc("/.well-known/jwks.json", jwksHandler)

	port := os.Getenv("PORT")

	fmt.Printf("Authorization Server running on http://localhost:%s/\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

// func revokeHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Revocation Endpoint"))
// }

// func introspectHandler(w http.ResponseWriter, r *http.Request) {
// 	response := map[string]interface{}{
// 		"active": true,
// 		"scope":  "read write",
// 		"exp":    1712345678,
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// func jwksHandler(w http.ResponseWriter, r *http.Request) {
// 	response := map[string]interface{}{
// 		"keys": []map[string]string{
// 			{
// 				"kty": "RSA",
// 				"kid": "example-key-id",
// 				"alg": "RS256",
// 				"n":   "example-modulus",
// 				"e":   "AQAB",
// 			},
// 		},
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
