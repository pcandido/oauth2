package token

import (
	"authorization-server/config"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

func Validate(token string) (*jwt.MapClaims, error) {
	secretKey := []byte(config.ACCESS_TOKEN_SECRET)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) { return secretKey, nil })
	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}
