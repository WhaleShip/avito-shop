package store

import (
	"context"
	"testing"

	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestUserFunctions(t *testing.T) {
	t.Run("GetUserByUsername: Успешное получение пользователя", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		// coins передаём как int64
		rows := pgxmock.NewRows([]string{"username", "password", "coins"}).
			AddRow("testuser", "hashed", int64(1000))

		query := "^SELECT username, password, coins FROM users WHERE username=\\$1$"
		mockConn.ExpectQuery(query).
			WithArgs("testuser").
			WillReturnRows(rows)

		user, err := GetUserByUsername(context.Background(), mockConn, "testuser")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if user.Username != "testuser" || user.Coins != int64(1000) {
			t.Errorf("неверный пользователь: %+v", user)
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("CreateUser: Успешное создание пользователя", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		query := "^INSERT INTO users\\(username, password, coins\\) VALUES\\(\\$1, \\$2, \\$3\\)$"
		mockConn.ExpectExec(query).
			WithArgs("newuser", "secret", int64(1000)).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = CreateUser(context.Background(), mockConn, "newuser", "secret")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("GetUserByUsernameTx: Успешное получение пользователя в транзакции", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		rows := pgxmock.NewRows([]string{"username", "password", "coins"}).
			AddRow("testuser", "hashed", int64(1000))

		query := "^SELECT username, password, coins FROM users WHERE username=\\$1 FOR UPDATE$"
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
		mockConn.ExpectCommit()
		if err = tx.Commit(context.Background()); err != nil {
			t.Errorf("ошибка коммита: %v", err)
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("UpdateUserCoinsTx: Успешное обновление количества монет", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		// Обновление по username
		mockConn.ExpectExec("^UPDATE users SET coins = \\$1 WHERE username = \\$2$").
			WithArgs(int64(500), "testuser").
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err = UpdateUserCoinsTx(context.Background(), tx, "testuser", 500)
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
