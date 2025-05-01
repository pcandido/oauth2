package token

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"maps"

	"github.com/golang-jwt/jwt/v4"
)

func Generate(claims map[string]any, expiration time.Duration) (string, error) {
	privateKey, err := getPrivateKey("private_key.pem")
	if err != nil {
		return "", err
	}

	tokenClaims := jwt.MapClaims{}
	maps.Copy(tokenClaims, claims)
	tokenClaims["exp"] = time.Now().Add(expiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenClaims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign the token: %w", err)
	}

	return signedToken, nil
}

func getPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privateKey, nil
}
