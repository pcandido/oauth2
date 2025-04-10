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
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/callback", handlers.LoginCallbackHandler)

	fmt.Printf("Authorization Client running on %s\n", config.Url(""))
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port()), nil)
}
