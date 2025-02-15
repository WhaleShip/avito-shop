package store

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestGetInventory(t *testing.T) {
	t.Run("Успешное получение списка товаров", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		var uid uint = 1
		rows := pgxmock.NewRows([]string{"id", "user_id", "item_name", "quantity"}).
			AddRow(uint(1), uid, "item1", 10).
			AddRow(uint(2), uid, "item2", 5)

		mockConn.ExpectQuery(`SELECT id, user_id, item_name, quantity FROM inventory_items WHERE user_id=\$1`).
			WithArgs(uid).
			WillReturnRows(rows)

		items, err := GetInventory(context.Background(), mockConn, uid)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("ожидалось 2 элемента, получено %d", len(items))
		}
	})

	t.Run("Ошибка выполнения запроса", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		defer CloseAndLogMock(mockConn)

		var uid uint = 1
		mockConn.ExpectQuery(`SELECT id, user_id, item_name, quantity FROM inventory_items WHERE user_id=\$1`).
			WithArgs(uid).
			WillReturnError(errors.New("query error"))

		items, err := GetInventory(context.Background(), mockConn, uid)
		if err == nil {
			t.Error("ожидалась ошибка, получено nil")
		}
		if items != nil {
			t.Errorf("ожидался nil для items, получено: %+v", items)
		}
	})
}

func TestUpsertInventoryItemTx(t *testing.T) {
	t.Run("INSERT: запись не найдена", func(t *testing.T) {
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
		mockConn.ExpectQuery(`SELECT id, quantity FROM inventory_items WHERE user_id=\$1 AND item_name=\$2 FOR UPDATE`).
			WithArgs(uid, "itemA").
			WillReturnError(pgx.ErrNoRows)

		mockConn.ExpectExec(`INSERT INTO inventory_items\(user_id, item_name, quantity\) VALUES\(\$1, \$2, 1\)`).
			WithArgs(uid, "itemA").
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = UpsertInventoryItemTx(context.Background(), tx, uid, "itemA")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		TxCommitAndLog(tx)
	})

	t.Run("UPDATE: запись найдена", func(t *testing.T) {
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
		rows := pgxmock.NewRows([]string{"id", "quantity"}).
			AddRow(uid, 2)
		mockConn.ExpectQuery(`SELECT id, quantity FROM inventory_items WHERE user_id=\$1 AND item_name=\$2 FOR UPDATE`).
			WithArgs(uid, "itemB").
			WillReturnRows(rows)

		mockConn.ExpectExec(`UPDATE inventory_items SET quantity = quantity \+ 1 WHERE id = \$1`).
			WithArgs(uid).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err = UpsertInventoryItemTx(context.Background(), tx, uid, "itemB")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		TxCommitAndLog(tx)
	})
}
