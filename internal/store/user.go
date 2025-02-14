package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func GetUserByUsername(ctx context.Context, db database.PgxIface, username string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(ctx,
		"SELECT id, username, password, coins, created_at FROM users WHERE username=$1",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.Coins, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(ctx context.Context, db database.PgxIface, username, password string) error {
	_, err := db.Exec(ctx,
		"INSERT INTO users(username, password, coins, created_at) VALUES($1, $2, $3, now())",
		username, password, 1000)
	return err
}

func GetUsernameByID(ctx context.Context, db database.PgxIface, userID uint) (string, error) {
	var username string
	err := db.QueryRow(ctx, "SELECT username FROM users WHERE id=$1", userID).Scan(&username)
	return username, err
}

func GetUserByUsernameTx(ctx context.Context, tx pgx.Tx, username string) (*models.User, error) {
	user := &models.User{}
	err := tx.QueryRow(ctx,
		"SELECT id, username, password, coins, created_at FROM users WHERE username=$1 FOR UPDATE",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.Coins, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserCoinsTx(ctx context.Context, tx pgx.Tx, userID uint, coins int) error {
	_, err := tx.Exec(ctx, "UPDATE users SET coins = $1 WHERE id = $2", coins, userID)
	return err
}
