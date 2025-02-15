package store

import (
	"context"
	"errors"
	"testing"
	"time"

	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestCreateCoinTransactionTx(t *testing.T) {
	t.Run("Успешное создание транзакции монет", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		var fromUser, toUser uint = 1, 2
		mockConn.ExpectExec(`INSERT INTO `+
			`coin_transactions\(from_user_id, to_user_id, amount, created_at\) `+
			`VALUES\(\$1, \$2, \$3, now\(\)\)`).
			WithArgs(fromUser, toUser, 50).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = CreateCoinTransactionTx(context.Background(), tx, fromUser, toUser, 50)
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		TxCommitAndLog(tx)
	})
}

func TestFinalizeTransaction(t *testing.T) {
	t.Run("Коммит транзакции при отсутствии ошибки", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		mockConn.ExpectCommit()
		FinalizeTransaction(nil, tx)
	})

	t.Run("Откат транзакции при наличии ошибки", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		mockConn.ExpectRollback()
		FinalizeTransaction(errors.New("some error"), tx)
	})
}

func TestGetCoinTransactions(t *testing.T) {
	t.Run("Сценарий 'sent'", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		var id, fromUser, toUser uint = 1, 1, 2
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "created_at"}).
			AddRow(id, fromUser, toUser, 50, now)

		query := `SELECT id, from_user_id, to_user_id, amount, created_at ` +
			`FROM coin_transactions ` +
			`WHERE from_user_id=\$1 ` +
			`ORDER BY created_at`
		mockConn.ExpectQuery(query).
			WithArgs(fromUser).
			WillReturnRows(rows)

		transactions, err := GetCoinTransactions(context.Background(), mockConn, fromUser, "sent")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("ожидалась 1 транзакция, получено %d", len(transactions))
		}
	})

	t.Run("Сценарий 'received'", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		var id, fromUser, toUser uint = 1, 2, 1
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "from_user_id", "to_user_id", "amount", "created_at"}).
			AddRow(id, fromUser, toUser, 30, now)

		query := `SELECT id, from_user_id, to_user_id, amount, created_at ` +
			`FROM coin_transactions ` +
			`WHERE to_user_id=\$1 ` +
			`ORDER BY created_at`
		mockConn.ExpectQuery(query).
			WithArgs(toUser).
			WillReturnRows(rows)

		transactions, err := GetCoinTransactions(context.Background(), mockConn, toUser, "received")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if len(transactions) != 1 {
			t.Errorf("ожидалась 1 транзакция, получено %d", len(transactions))
		}
	})

	t.Run("Неверное направление", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		transactions, err := GetCoinTransactions(context.Background(), mockConn, 1, "invalid")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if transactions != nil {
			t.Errorf("ожидался nil, получено: %+v", transactions)
		}
	})
}
