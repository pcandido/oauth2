package token

import (
	"authorization-server/config"
	"fmt"

	"time"

	"github.com/golang-jwt/jwt/v4"
)

func Generate(claims map[string]string, expiration time.Duration) (string, error) {
	secretKey := []byte(config.ACCESS_TOKEN_SECRET)

	tokenClaims := jwt.MapClaims{}
	for key, value := range claims {
		tokenClaims[key] = value
	}
	tokenClaims["exp"] = time.Now().Add(expiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign the token: %s", err.Error())
	}

	return signedToken, nil
}
