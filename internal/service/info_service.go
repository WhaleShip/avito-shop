package service

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/database/models"
	"github.com/whaleship/avito-shop/internal/dto"
	"github.com/whaleship/avito-shop/internal/store"
)

func GetUserInfo(ctx context.Context, db database.PgxIface, username string) (*dto.InfoResponse, error) {
	user, err := store.GetUserByUsername(ctx, db, username)
	if err != nil {
		return nil, err
	}

	invItems, err := store.GetInventory(ctx, db, user.ID)
	if err != nil {
		return nil, err
	}
	var inventoryResp []dto.InventoryItemResp
	for _, item := range invItems {
		inventoryResp = append(inventoryResp, dto.InventoryItemResp{
			Type:     item.ItemName,
			Quantity: item.Quantity,
		})
	}

	sentTx, err := store.GetCoinTransactions(ctx, db, user.ID, "sent")
	if err != nil {
		log.Println("error getting sent transactions: ", err)
		sentTx = []models.CoinTransaction{}
	}
	receivedTx, err := store.GetCoinTransactions(ctx, db, user.ID, "received")
	if err != nil {
		log.Println("error getting received transactions: ", err)
		receivedTx = []models.CoinTransaction{}
	}

	var sentResp []dto.SentTxResp
	for _, tx := range sentTx {
		toUser, err := store.GetUsernameByID(ctx, db, tx.ToUserID)
		if err != nil {
			toUser = ""
		}
		sentResp = append(sentResp, dto.SentTxResp{
			ToUser: toUser,
			Amount: tx.Amount,
		})
	}

	var receivedResp []dto.ReceivedTxResp
	for _, tx := range receivedTx {
		fromUser, err := store.GetUsernameByID(ctx, db, tx.FromUserID)
		if err != nil {
			fromUser = ""
		}
		receivedResp = append(receivedResp, dto.ReceivedTxResp{
			FromUser: fromUser,
			Amount:   tx.Amount,
		})
	}

	infoResp := &dto.InfoResponse{
		Coins:     user.Coins,
		Inventory: inventoryResp,
		CoinHistory: dto.CoinHistoryResp{
			Received: receivedResp,
			Sent:     sentResp,
		},
	}
	return infoResp, nil
}
