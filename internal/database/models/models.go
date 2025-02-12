package models

import "time"

type User struct {
	ID        uint      `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Coins     int       `db:"coins"`
	CreatedAt time.Time `db:"created_at"`
}

type CoinTransaction struct {
	ID         uint      `db:"id"`
	FromUserID uint      `db:"from_user_id"`
	ToUserID   uint      `db:"to_user_id"`
	Amount     int       `db:"amount"`
	CreatedAt  time.Time `db:"created_at"`
}

type InventoryItem struct {
	ID       uint   `db:"id"`
	UserID   uint   `db:"user_id"`
	ItemName string `db:"item_name"`
	Quantity int    `db:"quantity"`
}

type MerchItem struct {
	ID    uint   `db:"id"`
	Name  string `db:"name"`
	Price int    `db:"price"`
}
