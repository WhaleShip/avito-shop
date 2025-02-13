package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func GetMerchItemByNameTx(ctx context.Context, tx pgx.Tx, name string) (*models.MerchItem, error) {
	item := &models.MerchItem{}
	err := tx.QueryRow(ctx, "SELECT id, name, price FROM merch_items WHERE name=$1 FOR UPDATE", name).
		Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		return nil, err
	}
	return item, nil
}
