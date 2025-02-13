package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func GetInventory(ctx context.Context, db *pgx.Conn, userID uint) ([]models.InventoryItem, error) {
	rows, err := db.Query(ctx, "SELECT id, user_id, item_name, quantity FROM inventory_items WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.ItemName, &item.Quantity); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}
