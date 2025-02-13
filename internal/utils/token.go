package utils

import (
	"errors"
	"time"

	"github.com/whaleship/avito-shop/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to convert token to MapClaims")
	}
	mapClaims["username"] = username
	mapClaims["exp"] = time.Now().Add(72 * time.Hour).Unix()
	return token.SignedString([]byte(config.GetJWTSecret()))
}
