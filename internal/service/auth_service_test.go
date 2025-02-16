package service

import (
	"context"
	"errors"
	"testing"

	"bou.ke/monkey"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/database/models"
	"github.com/whaleship/avito-shop/internal/store"
	"github.com/whaleship/avito-shop/internal/utils"
)

func TestAuthenticateOrCreateUser(t *testing.T) {
	ctx := context.Background()
	username := "testuser"
	password := "password"
	hashedPassword := "hashedpassword"

	t.Run("пользователь не существует, успешное создание", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return nil, errors.New("not found")
			})
		defer patchGetUser.Unpatch()

		patchHash := monkey.Patch(utils.HashPassword, func(pwd string) (string, error) {
			return hashedPassword, nil
		})
		defer patchHash.Unpatch()

		mockConn.ExpectExec("INSERT INTO users").
			WithArgs(username, hashedPassword, pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = AuthenticateOrCreateUser(ctx, mockConn, username, password)
		if err != nil {
			t.Error("неожиданная ошибка: ", err)
		}
	})

	t.Run("пользователь существует и пароль верный", func(t *testing.T) {
		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return &models.User{Username: uname, Password: "hashedpassword"}, nil
			})
		defer patchGetUser.Unpatch()

		patchCheck := monkey.Patch(utils.CheckPassword, func(hashedPwd, pwd string) bool {
			return true
		})
		defer patchCheck.Unpatch()

		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		err = AuthenticateOrCreateUser(ctx, mockConn, username, password)
		if err != nil {
			t.Error("неожиданная ошибка: ", err)
		}
	})

	t.Run("пользователь существует, неверный пароль", func(t *testing.T) {
		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return &models.User{Username: uname, Password: "hashedpassword"}, nil
			})
		defer patchGetUser.Unpatch()

		patchCheck := monkey.Patch(utils.CheckPassword, func(hashedPwd, pwd string) bool {
			return false
		})
		defer patchCheck.Unpatch()

		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		err = AuthenticateOrCreateUser(ctx, mockConn, username, password)
		if err == nil {
			t.Error("ожидалась ошибка неверных учетных данных")
		}
		var invCredErr *utils.InvalidCredentialsError
		if !errors.As(err, &invCredErr) {
			t.Error("ожидалась ошибка типа InvalidCredentialsError, получено: ", err)
		}
	})
}
