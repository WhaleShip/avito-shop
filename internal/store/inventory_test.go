package store

import (
	"context"
	"errors"
	"testing"

	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestGetInventory(t *testing.T) {
	t.Run("Успешное получение списка товаров", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		// Используем username (а не uid)
		username := "testuser"
		rows := pgxmock.NewRows([]string{"id", "user_username", "item_name", "quantity"}).
			AddRow(int64(1), username, "item1", 10).
			AddRow(int64(2), username, "item2", 5)

		mockConn.ExpectQuery("^SELECT id, user_username, item_name, quantity FROM inventory_items WHERE user_username=\\$1$").
			WithArgs(username).
			WillReturnRows(rows)

		items, err := GetInventory(context.Background(), mockConn, username)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("ожидалось 2 элемента, получено %d", len(items))
		}

		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("Ошибка выполнения запроса", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		username := "testuser"
		mockConn.ExpectQuery("^SELECT id, user_username, item_name, quantity FROM inventory_items WHERE user_username=\\$1$").
			WithArgs(username).
			WillReturnError(errors.New("query error"))

		items, err := GetInventory(context.Background(), mockConn, username)
		if err == nil {
			t.Error("ожидалась ошибка, получено nil")
		}
		if items != nil {
			t.Errorf("ожидался nil для items, получено: %+v", items)
		}
		if err := mockConn.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestUpsertInventoryItemTx(t *testing.T) {
	t.Run("INSERT: запись не найдена", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		username := "testuser"
		itemName := "itemA"
		query := "^INSERT INTO inventory_items\\(user_username, item_name, quantity\\) " +
			"VALUES\\(\\$1, \\$2, 1\\) ON CONFLICT \\(user_username, item_name\\) " +
			"DO UPDATE SET quantity = inventory_items\\.quantity \\+ 1$"
		mockConn.ExpectExec(query).
			WithArgs(username, itemName).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = UpsertInventoryItemTx(context.Background(), tx, username, itemName)
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
