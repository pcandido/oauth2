package config

import (
	"fmt"
	"net/url"
	"os"
)

func AuthorizationServerUrl(path string, backend bool, queryParams *url.Values) string {
	host := os.Getenv("AUTHORIZATION_SERVER_HOST_FE")
	if backend {
		host = os.Getenv("AUTHORIZATION_SERVER_HOST_BE")
	}

	url := &url.URL{
		Scheme: "http",
		Host:   host,
		Path:   path,
	}

	if queryParams != nil {
		url.RawQuery = queryParams.Encode()
	}

	return url.String()
}

func Url(path string) string {
	return fmt.Sprintf("http://localhost:%s/%s", Port(), path)
}

func Port() string {
	return os.Getenv("PORT")
}
