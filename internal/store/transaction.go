package store

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func CreateCoinTransactionTx(ctx context.Context, tx pgx.Tx, fromUser, toUser string, amount int64) error {
	_, err := tx.Exec(ctx,
		"INSERT INTO coin_transactions(from_user, to_user, amount) VALUES($1, $2, $3)",
		fromUser, toUser, amount)
	return err
}

func FinalizeTransaction(err error, tx pgx.Tx) {
	if err != nil {
		if err = tx.Rollback(context.Background()); err != nil {
			log.Println("error during rollback db: ", err)
		}
	} else {
		if err = tx.Commit(context.Background()); err != nil {
			log.Println("error during commit: ", err)
		}
	}
}

func GetCoinTransactions(ctx context.Context, db database.PgxIface, username, direction string) ([]models.CoinTransaction, error) {
	var rows pgx.Rows
	var err error
	if direction == "sent" {
		query := "SELECT id, from_user, to_user, amount FROM coin_transactions WHERE from_user=$1 ORDER BY id"
		rows, err = db.Query(ctx, query, username)
	} else if direction == "received" {
		query := "SELECT id, from_user, to_user, amount FROM coin_transactions WHERE to_user=$1 ORDER BY id"
		rows, err = db.Query(ctx, query, username)
	} else {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.CoinTransaction
	for rows.Next() {
		var tx models.CoinTransaction
		if err := rows.Scan(&tx.ID, &tx.FromUser, &tx.ToUser, &tx.Amount); err != nil {
			continue
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
