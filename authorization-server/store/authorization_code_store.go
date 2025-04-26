package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type AuthorizationCode struct {
	Code        string
	UserId      string
	clientId    string
	redirectUri string
	Scope       string
	createdAt   int64
}

// In production, use a shared and persistent store (like a Redis) for authorization codes.
var codes = make(map[string]AuthorizationCode)

func GenerateCode(userId string, clientId string, redirectUri string, scope string) string {
	code := generateRandomCode(32)
	codes[code] = AuthorizationCode{
		Code:        code,
		UserId:      userId,
		clientId:    clientId,
		redirectUri: redirectUri,
		Scope:       scope,
		createdAt:   time.Now().Unix(),
	}
	return code
}

func GetCode(code string, clientId string, redirectUri string) (*AuthorizationCode, error) {
	authCode, exists := codes[code]
	if !exists {
		return nil, fmt.Errorf("authorization code not found")
	}

	if authCode.clientId != clientId {
		return nil, fmt.Errorf("client_id mismatch")
	}

	if authCode.redirectUri != redirectUri {
		return nil, fmt.Errorf("redirect_uri mismatch")
	}

	delete(codes, code)
	return &authCode, nil
}

func generateRandomCode(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
