package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/whaleship/avito-shop/internal/service"
	"github.com/whaleship/avito-shop/internal/store"
	"github.com/whaleship/avito-shop/internal/utils"
)

func BuyHandler(c *fiber.Ctx) error {
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

	merchName := c.Params("item")
	ctx := c.UserContext()

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Println("error beginning transaction:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка начала транзакции"})
	}

	defer func() {
		store.FinalizeTransaction(c.UserContext(), err, tx)
	}()

	err = service.ProcessBuyMerch(ctx, tx, username, merchName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": utils.CapitalizeFirst(err.Error())})
	}

	return c.SendStatus(fiber.StatusOK)
}
