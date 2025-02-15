package config

import "os"

var jwtSecret []byte

func InitConfig() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

func GetJWTSecret() *[]byte {
	return &jwtSecret
}
