package config

import (
	"os"
	"sync"
)

var (
	jwtSecret []byte
	once      sync.Once
)

func initJWTSecret() {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		secretStr = "secret"
	}
	jwtSecret = []byte(secretStr)
}

func GetJWTSecret() []byte {
	once.Do(initJWTSecret)
	return jwtSecret
}
