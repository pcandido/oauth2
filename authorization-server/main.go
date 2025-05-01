package main

import (
	"authorization-server/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/authorize", handlers.AuthorizeHandler)
	http.HandleFunc("/token", handlers.TokenHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	port := os.Getenv("PORT")

	log.Printf("Authorization Server running on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
