package main

import (
	"authorization-client/config"
	"authorization-client/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/auth_server/initiate", handlers.LoginInitiateHandler)
	http.HandleFunc("/auth_server/callback", handlers.LoginCallbackHandler)

	log.Printf("Authorization Client running on port %s", config.Port())
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port()), nil)
}
