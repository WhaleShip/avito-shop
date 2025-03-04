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

func TestProcessSendCoin(t *testing.T) {
	ctx := t.Context()
	senderUsername := "sender"
	receiverUsername := "receiver"
	amount := int64(50)

	t.Run("недостаточно средств", func(t *testing.T) {
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
				if uname == senderUsername {
					return &models.User{Username: senderUsername, Coins: 30}, nil
				}
				return &models.User{Username: receiverUsername, Coins: 100}, nil
			})
		defer patchGetUser.Unpatch()

		err = ProcessSendCoin(ctx, tx, senderUsername, receiverUsername, amount)
		if err == nil || err.Error() != "недостаточно средств" {
			t.Error("ожидалась ошибка 'недостаточно средств', получено: ", err)
		}
	})

	t.Run("отправка монет себе", func(t *testing.T) {
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
				return &models.User{Username: uname, Coins: 100}, nil
			})
		defer patchGetUser.Unpatch()

		err = ProcessSendCoin(ctx, tx, senderUsername, senderUsername, amount)
		if err == nil || err.Error() != "нельзя отправлять монеты себе" {
			t.Error("ожидалась ошибка 'нельзя отправлять монеты себе', получено: ", err)
		}
	})

	t.Run("ошибка обновления средств отправителя", func(t *testing.T) {
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
				if uname == senderUsername {
					return &models.User{Username: senderUsername, Coins: 100}, nil
				}
				return &models.User{Username: receiverUsername, Coins: 100}, nil
			})
		defer patchGetUser.Unpatch()

		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(50), senderUsername).
			WillReturnError(errors.New("ошибка обновления средств отправителя"))

		err = ProcessSendCoin(ctx, tx, senderUsername, receiverUsername, amount)
		if err == nil || err.Error() != "ошибка обновления средств отправителя" {
			t.Error("ожидалась ошибка 'ошибка обновления средств отправителя', получено: ", err)
		}
	})

	t.Run("ошибка обновления средств получателя", func(t *testing.T) {
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
				if uname == senderUsername {
					return &models.User{Username: senderUsername, Coins: 100}, nil
				}
				return &models.User{Username: receiverUsername, Coins: 100}, nil
			})
		defer patchGetUser.Unpatch()

		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(50), senderUsername).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(150), receiverUsername).
			WillReturnError(errors.New("ошибка обновления средств получателя"))

		err = ProcessSendCoin(ctx, tx, senderUsername, receiverUsername, amount)
		if err == nil || err.Error() != "ошибка обновления средств получателя" {
			t.Error("ожидалась ошибка обновления средств получателя, получено: ", err)
		}
	})

	t.Run("ошибка записи транзакции", func(t *testing.T) {
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
				if uname == senderUsername {
					return &models.User{Username: senderUsername, Coins: 100}, nil
				}
				return &models.User{Username: receiverUsername, Coins: 100}, nil
			})
		defer patchGetUser.Unpatch()

		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(50), senderUsername).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(150), receiverUsername).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockConn.ExpectExec("INSERT INTO coin_transactions").
			WithArgs(senderUsername, receiverUsername, amount).
			WillReturnError(errors.New("ошибка записи транзакции"))

		err = ProcessSendCoin(ctx, tx, senderUsername, receiverUsername, amount)
		if err == nil || err.Error() != "ошибка записи транзакции" {
			t.Error("ожидалась ошибка записи транзакции, получено: ", err)
		}
	})

	t.Run("успешная отправка монет", func(t *testing.T) {
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
				if uname == senderUsername {
					return &models.User{Username: senderUsername, Coins: 100}, nil
				}
				return &models.User{Username: receiverUsername, Coins: 100}, nil
			})
		defer patchGetUser.Unpatch()

		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(50), senderUsername).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockConn.ExpectExec("UPDATE users SET coins").
			WithArgs(int64(150), receiverUsername).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		mockConn.ExpectExec("INSERT INTO coin_transactions").
			WithArgs(senderUsername, receiverUsername, amount).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = ProcessSendCoin(ctx, tx, senderUsername, receiverUsername, amount)
		if err != nil {
			t.Error("неожиданная ошибка: ", err)
		}

		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error("не все ожидания были выполнены: ", err)
		}
	})
}
