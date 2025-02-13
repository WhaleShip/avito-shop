package handlers

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/database/models"
	"github.com/whaleship/avito-shop/internal/dto"
	"github.com/whaleship/avito-shop/internal/store"
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

	user, err := store.GetUserByUsername(context.Background(), db, username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"errors": "Пользователь не найден"})
	}

	invItems, err := store.GetInventory(context.Background(), db, user.ID)
	if err != nil {
		log.Println("error getting inventory: ", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"errors": "Ошибка при получении инвентаря"})
	}
	var inventoryResp []dto.InventoryItemResp
	for _, item := range invItems {
		inventoryResp = append(inventoryResp, dto.InventoryItemResp{
			Type:     item.ItemName,
			Quantity: item.Quantity,
		})
	}

	sentTx, err := store.GetCoinTransactions(context.Background(), db, user.ID, "sent")
	if err != nil {
		sentTx = []models.CoinTransaction{}
	}
	receivedTx, err := store.GetCoinTransactions(context.Background(), db, user.ID, "received")
	if err != nil {
		receivedTx = []models.CoinTransaction{}
	}

	var sentResp []dto.SentTxResp
	for _, tx := range sentTx {
		toUser, err := store.GetUsernameByID(context.Background(), db, tx.ToUserID)
		if err != nil {
			toUser = ""
		}
		sentResp = append(sentResp, dto.SentTxResp{
			ToUser: toUser,
			Amount: tx.Amount,
		})
	}

	var received []dto.ReceivedTxResp
	for _, tx := range receivedTx {
		fromUser, err := store.GetUsernameByID(context.Background(), db, tx.FromUserID)
		if err != nil {
			fromUser = ""
		}
		received = append(received, dto.ReceivedTxResp{
			FromUser: fromUser,
			Amount:   tx.Amount,
		})
	}

	resp := dto.InfoResponse{
		Coins:     user.Coins,
		Inventory: inventoryResp,
		CoinHistory: dto.CoinHistoryResp{
			Received: received,
			Sent:     sentResp,
		},
	}
	return c.JSON(resp)
}
