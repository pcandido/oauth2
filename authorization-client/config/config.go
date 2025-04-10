package config

import (
	"fmt"
	"os"
)

func AuthorizationServerUrl(path string, backend bool) string {
	if backend {
		return fmt.Sprintf("http://%s/%s", os.Getenv("AUTHORIZATION_SERVER_HOST_BE"), path)
	}

	return fmt.Sprintf("http://%s/%s", os.Getenv("AUTHORIZATION_SERVER_HOST_FE"), path)
}

func ClientID() string {
	return os.Getenv("CLIENT_ID")
}

func ClientSecret() string {
	return os.Getenv("CLIENT_SECRET")
}

func Url(path string) string {
	return fmt.Sprintf("http://localhost:%s/%s", Port(), path)
}

func Port() string {
	return os.Getenv("PORT")
}
