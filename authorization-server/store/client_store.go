package store

import "fmt"

type Client struct {
	ClientID             string
	ClientSecret         string
	RedirectURIs         []string
	GrantTypes           []string
	ResponseTypes        []string
	Scope                string
	AccessTokenLifetime  int
	RefreshTokenLifetime int
}

var clients = map[string]Client{
	"client_id_1": {
		ClientID:             "client_id_1",
		ClientSecret:         "client_secret_1",
		RedirectURIs:         []string{"http://localhost:8080/callback"},
		GrantTypes:           []string{"authorization_code"},
		ResponseTypes:        []string{"code"},
		Scope:                "read,write",
		AccessTokenLifetime:  3600,
		RefreshTokenLifetime: 7200,
	},
}

func GetClient(clientID string) (*Client, error) {
	client, exists := clients[clientID]
	if !exists {
		return nil, fmt.Errorf("client not found")
	}
	return &client, nil
}
