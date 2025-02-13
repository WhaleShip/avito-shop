package handlers

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/store"
	"github.com/whaleship/avito-shop/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

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

	var req SendCoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Неверный запрос"})
	}
	if req.ToUser == "" || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Неверные параметры запроса"})
	}

	tx, err := db.Begin(context.Background())
	if err != nil {
		log.Println("error beginning transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка начала транзакции"})
	}
	defer func() {
		store.FinalizeTransaction(err, tx)
	}()

	sender, err := store.GetUserByUsernameTx(context.Background(), tx, username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Отправитель не найден"})
	}
	if sender.Coins < req.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Недостаточно средств"})
	}

	receiver, err := store.GetUserByUsernameTx(context.Background(), tx, req.ToUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Получатель не найден"})
	}

	if receiver.ID == sender.ID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Нельзя отправлять монеты самому себе"})
	}

	err = store.UpdateUserCoinsTx(context.Background(), tx, sender.ID, sender.Coins-req.Amount)
	if err != nil {
		log.Println("error updating sender coins: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка обновления средств отправителя"})
	}
	err = store.UpdateUserCoinsTx(context.Background(), tx, receiver.ID, receiver.Coins+req.Amount)
	if err != nil {
		log.Println("error updating receiver coins: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка обновления средств получателя"})
	}

	err = store.CreateCoinTransactionTx(context.Background(), tx, sender.ID, receiver.ID, req.Amount)
	if err != nil {
		log.Println("error creating coin transaction: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка записи транзакции"})
	}

	return c.SendStatus(fiber.StatusOK)
}
