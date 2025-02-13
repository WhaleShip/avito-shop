package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUsername(c *fiber.Ctx) (string, error) {
	userToken, ok := c.Locals("user").(*jwt.Token)
	if !ok {
		return "", errors.New("faild to get username")
	}

	claimsMap, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("failed to convert token to MapClaims")
	}

	res, ok := claimsMap["username"].(string)
	if !ok {
		return "", errors.New("failed to convert username to string")
	}

	return res, nil
}

func JwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Неавторизован"})
}
