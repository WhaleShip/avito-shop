package handlers

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/service"
	"github.com/whaleship/avito-shop/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func InfoHandler(c *fiber.Ctx) error {
	db, err := utils.ExtractDB(c)
	if err != nil {
		log.Println("error extracting context:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Внутренняя ошибка сервера"})
	}

	username, err := utils.GetUsername(c)
	if err != nil {
		log.Println("error getting username: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Внутренняя ошибка сервера"})
	}

	infoResp, err := service.GetUserInfo(context.Background(), db, username)
	if err != nil {
		log.Println("Error getting user info: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка получения информации"})
	}
	return c.JSON(infoResp)
}
