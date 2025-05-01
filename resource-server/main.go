package main

import (
	"fmt"
	"log"
	"net/http"
	"resource-server/config"
	"resource-server/handlers"
)

func main() {
	http.HandleFunc("/resource", handlers.ResourceHandler)

	port := config.Port()

	log.Printf("Resource Server running on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
