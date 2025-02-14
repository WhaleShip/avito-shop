package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/whaleship/avito-shop/internal/store"
)

func ProcessBuyMerch(ctx context.Context, tx pgx.Tx, username, merchName string) error {
	user, err := store.GetUserByUsernameTx(ctx, tx, username)
	if err != nil {
		return errors.New("пользователь не найден")
	}
	merchItem, err := store.GetMerchItemByNameTx(ctx, tx, merchName)
	if err != nil {
		return errors.New("товар не найден")
	}
	if user.Coins < merchItem.Price {
		return errors.New("недостаточно средств для покупки")
	}
	if err = store.UpdateUserCoinsTx(ctx, tx, user.ID, user.Coins-merchItem.Price); err != nil {
		log.Println("error updating user coins:", err)
		return errors.New("ошибка обновления средств")
	}
	if err = store.UpsertInventoryItemTx(ctx, tx, user.ID, merchItem.Name); err != nil {
		log.Println("error updating inventory:", err)
		return errors.New("ошибка обновления инвентаря")
	}
	return nil
}
