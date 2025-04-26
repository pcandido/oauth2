package store

import (
	"crypto/rand"
	"encoding/hex"
)

// In production, use a shared and persistent store (like a Redis) for authorization codes.
var codes = make(map[string]string)

func GenerateCode(userId string) string {
	code := generateRandomCode(32)
	codes[code] = userId
	return code
}

func GetUserId(code string) (string, bool) {
	userId, exists := codes[code]
	if exists {
		delete(codes, code) // Remove the code after use
	}
	return userId, exists
}

func generateRandomCode(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
