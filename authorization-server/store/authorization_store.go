package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type Authorization struct {
	Code        string
	UserId      string
	ClientId    string
	RedirectUri string
	Scope       string
	CreatedAt   int64
}

var authorizations = make(map[string]Authorization)

func PushAuthorization(userId string, clientId string, redirectUri string, scope string) string {
	code := generateRandomCode(32)
	authorizations[code] = Authorization{
		Code:        code,
		UserId:      userId,
		ClientId:    clientId,
		RedirectUri: redirectUri,
		Scope:       scope,
		CreatedAt:   time.Now().Unix(),
	}
	return code
}

func PopAuthorization(code string) (*Authorization, error) {
	authorization, exists := authorizations[code]
	if !exists {
		return nil, fmt.Errorf("authorization code not found")
	}

	delete(authorizations, code)
	return &authorization, nil
}

func generateRandomCode(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
