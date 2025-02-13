package store

import (
	"context"

	"github.com/whaleship/avito-shop/internal/database/models"

	"github.com/jackc/pgx/v5"
)

func GetUserByUsername(ctx context.Context, db *pgx.Conn, username string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(ctx,
		"SELECT id, username, password, coins, created_at FROM users WHERE username=$1",
		username).Scan(&user.ID, &user.Username, &user.Password, &user.Coins, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(ctx context.Context, db *pgx.Conn, username, password string) error {
	_, err := db.Exec(ctx,
		"INSERT INTO users(username, password, coins, created_at) VALUES($1, $2, $3, now())",
		username, password, 1000)
	return err
}
