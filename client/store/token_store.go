package store

import (
	"authorization-client/utils"
	"fmt"
	"log"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

var tokens = make(map[string]Tokens)

func StoreToken(token Tokens) string {
	session := utils.GenerateRandomString(32)
	tokens[session] = token
	log.Println("Token stored for session:", session)

	return session
}

func GetToken(session string) (*Tokens, error) {
	if token, exists := tokens[session]; exists {
		return &token, nil
	}
	return nil, fmt.Errorf("token not found for session: %s", session)
}
