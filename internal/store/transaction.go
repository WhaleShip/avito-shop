package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

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
