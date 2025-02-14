package handlers

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/dto"
	"github.com/whaleship/avito-shop/internal/service"
	"github.com/whaleship/avito-shop/internal/store"
	"github.com/whaleship/avito-shop/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func SendCoinHandler(c *fiber.Ctx) error {
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

	var req dto.SendCoinRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("error parsing send coin request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Неверный запрос"})
	}
	if req.ToUser == "" || req.Amount <= 0 {
		log.Printf("invalid parameters in send coin request: %+v", req)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Неверные параметры запроса"})
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("error beginning transaction:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка начала транзакции"})
	}
	defer func() {
		store.FinalizeTransaction(err, tx)
	}()

	err = service.ProcessSendCoin(context.Background(), tx, username, req.ToUser, req.Amount)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": utils.CapitalizeFirst(err.Error())})
	}
	return c.SendStatus(fiber.StatusOK)
}
