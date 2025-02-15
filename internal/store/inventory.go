package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func GetInventory(ctx context.Context, db database.PgxIface, username string) ([]models.InventoryItem, error) {
	rows, err := db.Query(ctx,
		"SELECT id, user_username, item_name, quantity "+
			"FROM inventory_items "+
			"WHERE user_username=$1",
		username)
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
	_, err := tx.Exec(ctx,
		`INSERT INTO inventory_items(user_username, item_name, quantity)
		 VALUES($1, $2, 1)
		 ON CONFLICT (user_username, item_name)
		 DO UPDATE SET quantity = inventory_items.quantity + 1`,
		username, itemName)
	return err
}
