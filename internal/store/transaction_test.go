package store

import (
	"context"
	"errors"
	"testing"

	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestCreateCoinTransactionTx(t *testing.T) {
	t.Run("Успешное создание транзакции монет", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		fromUser, toUser := "user1", "user2"
		mockConn.ExpectExec("^INSERT INTO coin_transactions\\(from_user, to_user, amount\\) VALUES\\(\\$1, \\$2, \\$3\\)$").
			WithArgs(fromUser, toUser, int64(50)).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = CreateCoinTransactionTx(context.Background(), tx, fromUser, toUser, 50)
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		mockConn.ExpectCommit()
		if err = tx.Commit(context.Background()); err != nil {
			t.Errorf("ошибка коммита: %v", err)
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestFinalizeTransaction(t *testing.T) {
	t.Run("Коммит транзакции при отсутствии ошибки", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		mockConn.ExpectCommit()
		// Перед вызовом FinalizeTransaction не нужно самостоятельно вызывать Commit()
		FinalizeTransaction(nil, tx)
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Откат транзакции при наличии ошибки", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		mockConn.ExpectRollback()
		FinalizeTransaction(errors.New("some error"), tx)
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestGetCoinTransactions(t *testing.T) {
	t.Run("Сценарий 'sent'", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		// Используем int64 для id
		rows := pgxmock.NewRows([]string{"id", "from_user", "to_user", "amount"}).
			AddRow(int64(1), "user1", "user2", int64(50))

		query := "^SELECT id, from_user, to_user, amount FROM coin_transactions WHERE from_user=\\$1 ORDER BY id$"
		mockConn.ExpectQuery(query).
			WithArgs("user1").
			WillReturnRows(rows)

		transactions, err := GetCoinTransactions(context.Background(), mockConn, "user1", "sent")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("ожидалась 1 транзакция, получено %d", len(transactions))
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Сценарий 'received'", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		rows := pgxmock.NewRows([]string{"id", "from_user", "to_user", "amount"}).
			AddRow(int64(1), "user2", "user1", int64(30))

		query := "^SELECT id, from_user, to_user, amount FROM coin_transactions WHERE to_user=\\$1 ORDER BY id$"
		mockConn.ExpectQuery(query).
			WithArgs("user1").
			WillReturnRows(rows)

		transactions, err := GetCoinTransactions(context.Background(), mockConn, "user1", "received")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("ожидалась 1 транзакция, получено %d", len(transactions))
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Неверное направление", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		transactions, err := GetCoinTransactions(context.Background(), mockConn, "user1", "invalid")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if transactions != nil {
			t.Errorf("ожидался nil, получено: %+v", transactions)
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}
