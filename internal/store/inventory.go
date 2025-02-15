package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func GetInventory(ctx context.Context, db database.PgxIface, username string) ([]models.InventoryItem, error) {
	rows, err := db.Query(ctx, "SELECT id, user_username, item_name, quantity FROM inventory_items WHERE user_username=$1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		var id int64
		if err := rows.Scan(&id, &item.UserName, &item.ItemName, &item.Quantity); err != nil {
			continue
		}
		item.ID = uint(id)
		items = append(items, item)
	}
	return items, nil
}

func UpsertInventoryItemTx(ctx context.Context, tx pgx.Tx, username, itemName string) error {
	var id int64
	var quantity int
	err := tx.QueryRow(ctx,
		"SELECT id, quantity FROM inventory_items WHERE user_username=$1 AND item_name=$2 FOR UPDATE", username, itemName).
		Scan(&id, &quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, err = tx.Exec(ctx,
				"INSERT INTO inventory_items(user_username, item_name, quantity) VALUES($1, $2, 1)", username, itemName)
			return err
		}
		return err
	}
	_, err = tx.Exec(ctx,
		"UPDATE inventory_items SET quantity = quantity + 1 WHERE id = $1", id)
	return err
}
