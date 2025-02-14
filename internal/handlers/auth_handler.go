package handlers

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/dto"
	"github.com/whaleship/avito-shop/internal/service"
	"github.com/whaleship/avito-shop/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthHandler(c *fiber.Ctx) error {
	db, err := utils.ExtractDB(c)
	if err != nil {
		log.Println("error extracting context:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Внутренняя ошибка сервера"})
	}

	var req dto.AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Неверный запрос"})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Username и password обязательны"})
	}

	err = service.AuthenticateOrCreateUser(context.Background(), db, req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Неверный логин или пароль"})
		}
		log.Println("authentication error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка при аутентификации"})
	}

	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		log.Println("error generating token: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка при генерации токена"})
	}
	return c.JSON(dto.AuthResponse{Token: token})
}
