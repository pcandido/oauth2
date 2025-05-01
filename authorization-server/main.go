package main

import (
	"authorization-server/config"
	"authorization-server/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/authorize", handlers.AuthorizeHandler)
	http.HandleFunc("/token", handlers.TokenHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/consent", handlers.ConsentHandler)

	port := config.Port()

	log.Printf("Authorization Server running on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
