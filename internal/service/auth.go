package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/whaleship/avito-shop/internal/store"
	"github.com/whaleship/avito-shop/internal/utils"
)

func AuthenticateUser(ctx context.Context, db *pgx.Conn, username, password string) error {
	user, err := store.GetUserByUsername(ctx, db, username)
	if err != nil {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}
		err = store.CreateUser(ctx, db, username, hashedPassword)
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
	} else {
		if !utils.CheckPassword(user.Password, password) {
			return errors.New("invalid credentials")
		}
	}
	return nil
}
