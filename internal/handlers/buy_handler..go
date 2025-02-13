package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
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

	tx, err := db.Begin(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка начала транзакции"})
	}
	defer func() {
		store.FinalizeTransaction(err, tx)
	}()

	user, err := store.GetUserByUsernameTx(context.Background(), tx, username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Пользователь не найден"})
	}

	merchItem, err := store.GetMerchItemByNameTx(context.Background(), tx, merchName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Товар не найден"})
	}

	if user.Coins < merchItem.Price {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Недостаточно средств для покупки"})
	}

	err = store.UpdateUserCoinsTx(context.Background(), tx, user.ID, user.Coins-merchItem.Price)
	if err != nil {
		log.Println("error updating user coins: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка обновления средств"})
	}

	err = store.UpsertInventoryItemTx(context.Background(), tx, user.ID, merchItem.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка обновления инвентаря"})
	}

	return c.SendStatus(fiber.StatusOK)
}
