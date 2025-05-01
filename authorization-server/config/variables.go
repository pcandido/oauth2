package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found or failed to load: %v", err)
	}
	log.Printf("Environment variables loaded from .env file")
}

func Port() string {
	return os.Getenv("PORT")
}
