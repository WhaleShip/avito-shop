package service

import (
	"context"
	"errors"
	"testing"

	"bou.ke/monkey"
	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/whaleship/avito-shop/internal/database/models"
	"github.com/whaleship/avito-shop/internal/store"
)

func TestProcessBuyMerch(t *testing.T) {
	ctx := context.Background()
	username := "buyer"
	merchName := "merchA"

	t.Run("пользователь не найден", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(ctx)
		if err != nil {
			t.Fatal("ошибка начала транзакции: ", err)
		}
		patchGetUser := monkey.Patch(store.GetUserByUsernameTx,
			func(ctx context.Context, tx pgx.Tx, uname string) (*models.User, error) {
				return nil, errors.New("пользователь не найден")
			})
		defer patchGetUser.Unpatch()

		err = ProcessBuyMerch(ctx, tx, username, merchName)
		if err == nil || err.Error() != "пользователь не найден" {
			t.Error("ожидалась ошибка 'пользователь не найден', получено: ", err)
		}

	})

	t.Run("товар не найден", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(ctx)
		if err != nil {
			t.Fatal("ошибка начала транзакции: ", err)
		}
		patchGetUser := monkey.Patch(store.GetUserByUsernameTx,
			func(ctx context.Context, tx pgx.Tx, uname string) (*models.User, error) {
				return &models.User{Username: uname, Coins: 200}, nil
			})
		defer patchGetUser.Unpatch()

		patchGetMerch := monkey.Patch(store.GetMerchItemByNameTx,
			func(ctx context.Context, tx pgx.Tx, mName string) (*models.MerchItem, error) {
				return nil, errors.New("товар не найден")
			})
		defer patchGetMerch.Unpatch()

		err = ProcessBuyMerch(ctx, tx, username, merchName)
		if err == nil || err.Error() != "товар не найден" {
			t.Error("ожидалась ошибка 'товар не найден', получено: ", err)
		}

	})

	t.Run("недостаточно средств для покупки", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(ctx)
		if err != nil {
			t.Fatal("ошибка начала транзакции: ", err)
		}
		patchGetUser := monkey.Patch(store.GetUserByUsernameTx,
			func(ctx context.Context, tx pgx.Tx, uname string) (*models.User, error) {
				return &models.User{Username: uname, Coins: 50}, nil
			})
		defer patchGetUser.Unpatch()

		patchGetMerch := monkey.Patch(store.GetMerchItemByNameTx,
			func(ctx context.Context, tx pgx.Tx, mName string) (*models.MerchItem, error) {
				return &models.MerchItem{Name: mName, Price: 100}, nil
			})
		defer patchGetMerch.Unpatch()

		err = ProcessBuyMerch(ctx, tx, username, merchName)
		if err == nil || err.Error() != "недостаточно средств для покупки" {
			t.Error("ожидалась ошибка 'недостаточно средств для покупки', получено: ", err)
		}

	})

	t.Run("ошибка обновления средств", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(ctx)
		if err != nil {
			t.Fatal("ошибка начала транзакции: ", err)
		}
		patchGetUser := monkey.Patch(store.GetUserByUsernameTx,
			func(ctx context.Context, tx pgx.Tx, uname string) (*models.User, error) {
				return &models.User{Username: uname, Coins: 200}, nil
			})
		defer patchGetUser.Unpatch()

		patchGetMerch := monkey.Patch(store.GetMerchItemByNameTx,
			func(ctx context.Context, tx pgx.Tx, mName string) (*models.MerchItem, error) {
				return &models.MerchItem{Name: mName, Price: 100}, nil
			})
		defer patchGetMerch.Unpatch()

		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(100), username).
			WillReturnError(errors.New("ошибка обновления средств"))

		err = ProcessBuyMerch(ctx, tx, username, merchName)
		if err == nil || err.Error() != "ошибка обновления средств" {
			t.Error("ожидалась ошибка 'ошибка обновления средств', получено: ", err)
		}

	})

	t.Run("ошибка обновления инвентаря", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(ctx)
		if err != nil {
			t.Fatal("ошибка начала транзакции: ", err)
		}
		patchGetUser := monkey.Patch(store.GetUserByUsernameTx,
			func(ctx context.Context, tx pgx.Tx, uname string) (*models.User, error) {
				return &models.User{Username: uname, Coins: 200}, nil
			})
		defer patchGetUser.Unpatch()

		patchGetMerch := monkey.Patch(store.GetMerchItemByNameTx,
			func(ctx context.Context, tx pgx.Tx, mName string) (*models.MerchItem, error) {
				return &models.MerchItem{Name: mName, Price: 100}, nil
			})
		defer patchGetMerch.Unpatch()

		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(100), username).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockConn.ExpectExec("INSERT INTO inventory").
			WithArgs(username, merchName).
			WillReturnError(errors.New("ошибка обновления инвентаря"))

		err = ProcessBuyMerch(ctx, tx, username, merchName)
		if err == nil || err.Error() != "ошибка обновления инвентаря" {
			t.Error("ожидалась ошибка 'ошибка обновления инвентаря', получено: ", err)
		}

	})

	t.Run("успешная покупка товара", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(ctx)
		if err != nil {
			t.Fatal("ошибка начала транзакции: ", err)
		}
		patchGetUser := monkey.Patch(store.GetUserByUsernameTx,
			func(ctx context.Context, tx pgx.Tx, uname string) (*models.User, error) {
				return &models.User{Username: uname, Coins: 200}, nil
			})
		defer patchGetUser.Unpatch()

		patchGetMerch := monkey.Patch(store.GetMerchItemByNameTx,
			func(ctx context.Context, tx pgx.Tx, mName string) (*models.MerchItem, error) {
				return &models.MerchItem{Name: mName, Price: 100}, nil
			})
		defer patchGetMerch.Unpatch()

		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(100), username).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockConn.ExpectExec("INSERT INTO inventory").
			WithArgs(username, merchName).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = ProcessBuyMerch(ctx, tx, username, merchName)
		if err != nil {
			t.Error("неожиданная ошибка: ", err)
		}

		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error("не все ожидания были выполнены: ", err)
		}
	})
}
