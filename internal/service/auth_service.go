package service

import (
	"context"
	"fmt"

	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/store"
	"github.com/whaleship/avito-shop/internal/utils"
)

type InvalidCredentialsError struct{}

func (e *InvalidCredentialsError) Error() string {
	return "invalid credentials"
}

func AuthenticateOrCreateUser(ctx context.Context, db database.PgxIface, username, password string) error {
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
			return &utils.InvalidCredentialsError{}
		}
	}
	return nil
}
