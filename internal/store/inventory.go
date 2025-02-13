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

func UpsertInventoryItemTx(ctx context.Context, tx pgx.Tx, userID uint, itemName string) error {
	var id uint
	var quantity int
	err := tx.QueryRow(ctx,
		"SELECT id, quantity FROM inventory_items WHERE user_id=$1 AND item_name=$2 FOR UPDATE", userID, itemName).
		Scan(&id, &quantity)
	if err != nil {
		if err == pgx.ErrNoRows {
			_, err = tx.Exec(ctx,
				"INSERT INTO inventory_items(user_id, item_name, quantity) VALUES($1, $2, 1)", userID, itemName)
			return err
		}
		return err
	}
	_, err = tx.Exec(ctx,
		"UPDATE inventory_items SET quantity = quantity + 1 WHERE id = $1", id)
	return err
}
