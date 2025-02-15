package store

import (
	"context"
	"testing"

	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestGetMerchItemByNameTx(t *testing.T) {
	t.Run("Успешное получение мерча", func(t *testing.T) {
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatalf("ошибка создания mock соединения: %v", err)
		}
		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		rows := pgxmock.NewRows([]string{"id", "name", "price"}).
			AddRow(int64(1), "merch1", int64(100))
		mockConn.ExpectQuery("^SELECT id, name, price FROM merch_items WHERE name=\\$1 FOR UPDATE$").
			WithArgs("merch1").
			WillReturnRows(rows)

		item, err := GetMerchItemByNameTx(context.Background(), tx, "merch1")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if item == nil || item.Name != "merch1" || item.Price != int64(100) {
			t.Errorf("неверный объект мерча: %+v", item)
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
