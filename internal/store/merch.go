package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func GetMerchItemByNameTx(ctx context.Context, tx pgx.Tx, name string) (*models.MerchItem, error) {
	item := &models.MerchItem{}
	var id int64
	err := tx.QueryRow(ctx, "SELECT id, name, price FROM merch_items WHERE name=$1 FOR UPDATE", name).
		Scan(&id, &item.Name, &item.Price)
	if err != nil {
		return nil, err
	}
	item.ID = uint(id)
	return item, nil
}
