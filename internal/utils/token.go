package utils

import (
	"fmt"
	"time"

	"github.com/whaleship/avito-shop/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string) (*string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	claims := TokenClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}
	return &signedToken, nil
}
