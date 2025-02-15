package service

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/whaleship/avito-shop/internal/store"
)

func ProcessSendCoin(ctx context.Context, tx pgx.Tx, senderUsername, receiverUsername string, amount int64) error {
	sender, err := store.GetUserByUsernameTx(ctx, tx, senderUsername)
	if err != nil {
		return errors.New("отправитель не найден")
	}
	if sender.Coins < amount {
		return errors.New("недостаточно средств")
	}
	receiver, err := store.GetUserByUsernameTx(ctx, tx, receiverUsername)
	if err != nil {
		return errors.New("получатель не найден")
	}
	if receiver.Username == sender.Username {
		return errors.New("нельзя отправлять монеты себе")
	}
	if err = store.UpdateUserCoinsTx(ctx, tx, senderUsername, sender.Coins-amount); err != nil {
		log.Println("error updating sender coins:", err)
		return errors.New("ошибка обновления средств отправителя")
	}
	if err = store.UpdateUserCoinsTx(ctx, tx, receiverUsername, receiver.Coins+amount); err != nil {
		log.Println("error updating receiver coins:", err)
		return errors.New("ошибка обновления средств получателя")
	}
	if err = store.CreateCoinTransactionTx(ctx, tx, senderUsername, receiverUsername, amount); err != nil {
		log.Println("error creating coin transaction:", err)
		return errors.New("ошибка записи транзакции")
	}
	return nil
}
