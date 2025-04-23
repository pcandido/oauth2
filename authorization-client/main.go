package main

import (
	"authorization-client/config"
	"authorization-client/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/login/initiate", handlers.LoginInitiateHandler)
	http.HandleFunc("/login/callback", handlers.LoginCallbackHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	fmt.Printf("Authorization Client running on %s\n", config.Url(""))
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port()), nil)
}
