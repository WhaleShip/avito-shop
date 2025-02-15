package store

import (
	"context"
	"testing"
	"time"

	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestUserFunctions(t *testing.T) {
	t.Run("GetUserByUsername: Успешное получение пользователя", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		var uid uint = 1
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "username", "password", "coins", "created_at"}).
			AddRow(uid, "testuser", "hashed", 1000, now)

		query := `SELECT id, username, password, coins, created_at FROM users WHERE username=\$1`
		mockConn.ExpectQuery(query).
			WithArgs("testuser").
			WillReturnRows(rows)

		user, err := GetUserByUsername(context.Background(), mockConn, "testuser")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if user.Username != "testuser" || user.Coins != 1000 {
			t.Errorf("неверный пользователь: %+v", user)
		}
	})

	t.Run("CreateUser: Успешное создание пользователя", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		query := `INSERT INTO users\(username, password, coins, created_at\) VALUES\(\$1, \$2, \$3, now\(\)\)`
		mockConn.ExpectExec(query).
			WithArgs("newuser", "secret", 1000).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = CreateUser(context.Background(), mockConn, "newuser", "secret")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
	})

	t.Run("GetUsernameByID: Успешное получение username по ID", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		query := `SELECT username FROM users WHERE id=\$1`
		mockConn.ExpectQuery(query).
			WithArgs(uint(1)).
			WillReturnRows(pgxmock.NewRows([]string{"username"}).AddRow("testuser"))

		username, err := GetUsernameByID(context.Background(), mockConn, 1)
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if username != "testuser" {
			t.Errorf("ожидалось 'testuser', получено '%s'", username)
		}
	})

	t.Run("GetUserByUsernameTx: Успешное получение пользователя в транзакции", func(t *testing.T) {
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

		var uid uint = 1
		now := time.Now()
		rows := pgxmock.NewRows([]string{"id", "username", "password", "coins", "created_at"}).
			AddRow(uid, "testuser", "hashed", 1000, now)

		query := `SELECT id, username, password, coins, created_at FROM users WHERE username=\$1 FOR UPDATE`
		mockConn.ExpectQuery(query).
			WithArgs("testuser").
			WillReturnRows(rows)

		user, err := GetUserByUsernameTx(context.Background(), tx, "testuser")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if user.Username != "testuser" {
			t.Errorf("ожидалось 'testuser', получено '%s'", user.Username)
		}
		TxCommitAndLog(tx)
	})

	t.Run("UpdateUserCoinsTx: Успешное обновление количества монет", func(t *testing.T) {
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

		mockConn.ExpectExec(`UPDATE users SET coins = \$1 WHERE id = \$2`).
			WithArgs(500, uint(1)).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err = UpdateUserCoinsTx(context.Background(), tx, 1, 500)
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}

		TxCommitAndLog(tx)
	})
}
