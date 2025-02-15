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
		defer CloseAndLogMock(mockConn)

		mockConn.ExpectBegin()
		tx, err := mockConn.Begin(context.Background())
		if err != nil {
			t.Fatalf("ошибка начала транзакции: %v", err)
		}

		var merchID uint = 1
		rows := pgxmock.NewRows([]string{"id", "name", "price"}).
			AddRow(merchID, "merch1", 100)
		mockConn.ExpectQuery(`SELECT id, name, price FROM merch_items WHERE name=\$1 FOR UPDATE`).
			WithArgs("merch1").
			WillReturnRows(rows)

		item, err := GetMerchItemByNameTx(context.Background(), tx, "merch1")
		if err != nil {
			t.Errorf("неожиданная ошибка: %v", err)
		}
		if item == nil || item.Name != "merch1" || item.Price != 100 {
			t.Errorf("неверный объект мерча: %+v", item)
		}
		TxCommitAndLog(tx)
	})
}
