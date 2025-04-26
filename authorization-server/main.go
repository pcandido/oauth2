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

	port := os.Getenv("PORT")

	fmt.Printf("Authorization Server running on http://localhost:%s/\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
