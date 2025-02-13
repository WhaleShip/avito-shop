package store

import (
	"context"
	"log"

	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func CreateCoinTransactionTx(ctx context.Context, tx pgx.Tx, fromUserID, toUserID uint, amount int) error {
	_, err := tx.Exec(ctx,
		"INSERT INTO coin_transactions(from_user_id, to_user_id, amount, created_at) VALUES($1, $2, $3, now())",
		fromUserID, toUserID, amount)
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

func GetCoinTransactions(
	ctx context.Context,
	db *pgx.Conn, userID uint,
	direction string,
) ([]models.CoinTransaction, error) {
	var rows pgx.Rows
	var err error
	if direction == "sent" {
		query := "SELECT id, from_user_id, to_user_id, amount, created_at " +
			"FROM coin_transactions " +
			"WHERE from_user_id=$1 " +
			"ORDER BY created_at"
		rows, err = db.Query(ctx, query, userID)
	} else if direction == "received" {
		query := "SELECT id, from_user_id, to_user_id, amount, created_at " +
			"FROM coin_transactions " +
			"WHERE to_user_id=$1 " +
			"ORDER BY created_at"
		rows, err = db.Query(ctx, query, userID)
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
		if err := rows.Scan(&tx.ID, &tx.FromUserID, &tx.ToUserID, &tx.Amount, &tx.CreatedAt); err != nil {
			continue
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
