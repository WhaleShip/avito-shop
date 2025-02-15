package utils

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	username := "testuser"
	tokenStr, err := GenerateToken(username)
	if err != nil {
		t.Errorf("неожиданная ошибка при генерации токена: %v", err)
	}
	token, err := jwt.Parse(*tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		t.Errorf("неожиданная ошибка при разборе токена: %v", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["username"] != username {
			t.Errorf("ожидалось username '%s', получено '%s'", username, claims["username"])
		}
		// Проверяем наличие exp и что он больше текущего времени
		if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
			t.Error("ожидалось, что exp присутствует и больше текущего времени")
		}
	} else {
		t.Error("не удалось получить claims из токена")
	}
}
