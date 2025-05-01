package store

import "fmt"

type Client struct {
	ClientID             string
	ClientSecret         string
	RedirectURIs         []string
	GrantTypes           []string
	ResponseTypes        []string
	Scopes               []string
	AccessTokenLifetime  int
	RefreshTokenLifetime int
}

var clients = []Client{
	{
		ClientID:             "client_id",
		ClientSecret:         "client_secret",
		RedirectURIs:         []string{"http://localhost:8080/login/callback"},
		GrantTypes:           []string{"authorization_code"},
		ResponseTypes:        []string{"code"},
		Scopes:               []string{"read", "write"},
		AccessTokenLifetime:  3600,
		RefreshTokenLifetime: 7200,
	},
}

func GetClient(clientID string) (*Client, error) {
	for _, client := range clients {
		if client.ClientID == clientID {
			return &client, nil
		}
	}
	return nil, fmt.Errorf("client not found")
}
