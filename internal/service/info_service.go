package service

import (
	"context"
	"log"
	"sync"

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

	var (
		invItems                     []models.InventoryItem
		sentTx                       []models.CoinTransaction
		receivedTx                   []models.CoinTransaction
		wg                           sync.WaitGroup
		errInv, errSent, errReceived error
	)

	wg.Add(3)
	go func() {
		defer wg.Done()
		invItems, errInv = store.GetInventory(ctx, db, user.Username)
	}()
	go func() {
		defer wg.Done()
		sentTx, errSent = store.GetCoinTransactions(ctx, db, user.Username, "sent")
	}()
	go func() {
		defer wg.Done()
		receivedTx, errReceived = store.GetCoinTransactions(ctx, db, user.Username, "received")
	}()
	wg.Wait()

	if errInv != nil {
		return nil, errInv
	}
	if errSent != nil {
		log.Println("error getting sent transactions: ", errSent)
		sentTx = []models.CoinTransaction{}
	}
	if errReceived != nil {
		log.Println("error getting received transactions: ", errReceived)
		receivedTx = []models.CoinTransaction{}
	}

	var inventoryResp []dto.InventoryItemResp
	for _, item := range invItems {
		inventoryResp = append(inventoryResp, dto.InventoryItemResp{
			Type:     item.ItemName,
			Quantity: item.Quantity,
		})
	}

	var sentResp []dto.SentTxResp
	for _, tx := range sentTx {
		toUser := tx.ToUser
		sentResp = append(sentResp, dto.SentTxResp{
			ToUser: toUser,
			Amount: tx.Amount,
		})
	}

	var receivedResp []dto.ReceivedTxResp
	for _, tx := range receivedTx {
		fromUser := tx.FromUser
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
