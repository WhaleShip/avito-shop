package handlers

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/service"
	"github.com/whaleship/avito-shop/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func AuthHandler(c *fiber.Ctx) error {
	dbConn := c.Locals("db")
	if dbConn == nil {
		log.Println("DB connection not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Внутренняя ошибка сервера"})
	}
	db, ok := dbConn.(*pgx.Conn)
	if !ok {
		log.Println("DB connection type assertion failed")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Внутренняя ошибка сервера"})
	}

	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Неверный запрос"})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Username и password обязательны"})
	}

	err := service.AuthenticateUser(context.Background(), db, req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": "Неверный логин или пароль"})
		}
		log.Println("Authentication error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка при аутентификации"})
	}

	token, err := utils.GenerateToken(req.Username)
	if err != nil {
		log.Println("Error generating token: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка при генерации токена"})
	}
	return c.JSON(AuthResponse{Token: token})
}
