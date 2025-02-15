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
		"SELECT username, password, coins FROM users WHERE username=$1",
		username).Scan(&user.Username, &user.Password, &user.Coins)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(ctx context.Context, db database.PgxIface, username, password string) error {
	_, err := db.Exec(ctx,
		"INSERT INTO users(username, password, coins) VALUES($1, $2, $3)",
		username, password, 1000)
	return err
}

func GetUserByUsernameTx(ctx context.Context, tx pgx.Tx, username string) (*models.User, error) {
	user := &models.User{}
	err := tx.QueryRow(ctx,
		"SELECT username, password, coins FROM users WHERE username=$1 FOR UPDATE",
		username).Scan(&user.Username, &user.Password, &user.Coins)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserCoinsTx(ctx context.Context, tx pgx.Tx, username string, coins int64) error {
	_, err := tx.Exec(ctx, "UPDATE users SET coins = $1 WHERE username = $2", coins, username)
	return err
}
