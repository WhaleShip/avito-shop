package service

import (
	"context"
	"errors"
	"testing"

	"bou.ke/monkey"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/whaleship/avito-shop/internal/database"
	"github.com/whaleship/avito-shop/internal/database/models"
	"github.com/whaleship/avito-shop/internal/store"
)

func TestGetUserInfo(t *testing.T) {
	ctx := context.Background()

	t.Run("Пользователь не найден", func(t *testing.T) {
		username := "nonexistent"
		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return nil, errors.New("пользователь не найден")
			})
		defer patchGetUser.Unpatch()

		res, err := GetUserInfo(ctx, mockConn, username)
		if err == nil || err.Error() != "пользователь не найден" {
			t.Error("ожидалась ошибка 'пользователь не найден', получено: ", err)
		}
		if res != nil {
			t.Error("ожидался nil результат, получено: ", res)
		}
	})

	t.Run("Ошибка получения инвентаря", func(t *testing.T) {
		username := "testuser"
		dummyUser := &models.User{Username: username, Coins: 100}

		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return dummyUser, nil
			})
		defer patchGetUser.Unpatch()

		patchInventory := monkey.Patch(store.GetInventory,
			func(ctx context.Context, db database.PgxIface, uname string) ([]models.InventoryItem, error) {
				return nil, errors.New("ошибка получения инвентаря")
			})
		defer patchInventory.Unpatch()

		patchTx := monkey.Patch(store.GetCoinTransactions,
			func(ctx context.Context, db database.PgxIface, uname, txType string) ([]models.CoinTransaction, error) {
				return []models.CoinTransaction{}, nil
			})
		defer patchTx.Unpatch()

		res, err := GetUserInfo(ctx, mockConn, username)
		if err == nil || err.Error() != "ошибка получения инвентаря" {
			t.Error("ожидалась ошибка 'ошибка получения инвентаря', получено: ", err)
		}
		if res != nil {
			t.Error("ожидался nil результат, получено: ", res)
		}
	})

	t.Run("Ошибка получения отправленных транзакций", func(t *testing.T) {
		username := "testuser"
		dummyUser := &models.User{Username: username, Coins: 100}

		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return dummyUser, nil
			})
		defer patchGetUser.Unpatch()

		inventory := []models.InventoryItem{
			{ItemName: "item1", Quantity: 2},
		}
		patchInventory := monkey.Patch(store.GetInventory,
			func(ctx context.Context, db database.PgxIface, uname string) ([]models.InventoryItem, error) {
				return inventory, nil
			})
		defer patchInventory.Unpatch()

		patchTx := monkey.Patch(store.GetCoinTransactions,
			func(ctx context.Context, db database.PgxIface, uname, txType string) ([]models.CoinTransaction, error) {
				if txType == "sent" {
					return nil, errors.New("ошибка получения отправленных транзакций")
				}
				return []models.CoinTransaction{
					{FromUser: "userA", Amount: 20},
				}, nil
			})
		defer patchTx.Unpatch()

		res, err := GetUserInfo(ctx, mockConn, username)
		if err != nil {
			t.Error("неожиданная ошибка: ", err)
		}
		if len(res.CoinHistory.Sent) != 0 {
			t.Error("ожидалось, что отправленные транзакции будут пустыми, получено: ", res.CoinHistory.Sent)
		}
		if len(res.CoinHistory.Received) != 1 {
			t.Error("ожидалось, что полученные транзакции будут содержать 1 элемент, получено: ", res.CoinHistory.Received)
		}
	})

	t.Run("Ошибка получения полученных транзакций", func(t *testing.T) {
		username := "testuser"
		dummyUser := &models.User{Username: username, Coins: 100}

		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}
		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return dummyUser, nil
			})
		defer patchGetUser.Unpatch()

		inventory := []models.InventoryItem{
			{ItemName: "item1", Quantity: 2},
		}
		patchInventory := monkey.Patch(store.GetInventory,
			func(ctx context.Context, db database.PgxIface, uname string) ([]models.InventoryItem, error) {
				return inventory, nil
			})
		defer patchInventory.Unpatch()

		patchTx := monkey.Patch(store.GetCoinTransactions,
			func(ctx context.Context, db database.PgxIface, uname, txType string) ([]models.CoinTransaction, error) {
				if txType == "received" {
					return nil, errors.New("ошибка получения полученных транзакций")
				}
				return []models.CoinTransaction{
					{ToUser: "userB", Amount: 30},
				}, nil
			})
		defer patchTx.Unpatch()

		res, err := GetUserInfo(ctx, mockConn, username)
		if err != nil {
			t.Error("неожиданная ошибка: ", err)
		}
		if len(res.CoinHistory.Received) != 0 {
			t.Error("ожидалось, что полученные транзакции будут пустыми, получено: ", res.CoinHistory.Received)
		}
		if len(res.CoinHistory.Sent) != 1 {
			t.Error("ожидалось, что отправленные транзакции будут содержать 1 элемент, получено: ", res.CoinHistory.Sent)
		}
	})

	t.Run("Успешное получение информации о пользователе", func(t *testing.T) {
		username := "testuser"
		dummyUser := &models.User{Username: username, Coins: 150}

		mockConn, err := pgxmock.NewConn()
		if err != nil {
			t.Fatal("ошибка создания mock соединения: ", err)
		}

		patchGetUser := monkey.Patch(store.GetUserByUsername,
			func(ctx context.Context, db database.PgxIface, uname string) (*models.User, error) {
				return dummyUser, nil
			})
		defer patchGetUser.Unpatch()

		inventory := []models.InventoryItem{
			{ItemName: "item1", Quantity: 1},
			{ItemName: "item2", Quantity: 2},
		}
		patchInventory := monkey.Patch(store.GetInventory,
			func(ctx context.Context, db database.PgxIface, uname string) ([]models.InventoryItem, error) {
				return inventory, nil
			})
		defer patchInventory.Unpatch()

		patchTx := monkey.Patch(store.GetCoinTransactions,
			func(ctx context.Context, db database.PgxIface, uname, txType string) ([]models.CoinTransaction, error) {
				if txType == "sent" {
					return []models.CoinTransaction{
						{ToUser: "userA", Amount: 30},
					}, nil
				}
				if txType == "received" {
					return []models.CoinTransaction{
						{FromUser: "userB", Amount: 40},
					}, nil
				}
				return nil, nil
			})
		defer patchTx.Unpatch()

		res, err := GetUserInfo(ctx, mockConn, username)
		if err != nil {
			t.Error("неожиданная ошибка: ", err)
		}
		if res.Coins != 150 {
			t.Error("ожидалось, что у пользователя 150 монет, получено: ", res.Coins)
		}
		if len(res.Inventory) != 2 {
			t.Error("ожидалось 2 элемента инвентаря, получено: ", len(res.Inventory))
		}
		if len(res.CoinHistory.Sent) != 1 {
			t.Error("ожидалось 1 отправленная транзакция, получено: ", len(res.CoinHistory.Sent))
		}
		if len(res.CoinHistory.Received) != 1 {
			t.Error("ожидалось 1 полученная транзакция, получено: ", len(res.CoinHistory.Received))
		}
		if res.CoinHistory.Sent[0].ToUser != "userA" || res.CoinHistory.Sent[0].Amount != 30 {
			t.Error("неожиданные данные отправленной транзакции: ", res.CoinHistory.Sent[0])
		}
		if res.CoinHistory.Received[0].FromUser != "userB" || res.CoinHistory.Received[0].Amount != 40 {
			t.Error("неожиданные данные полученной транзакции: ", res.CoinHistory.Received[0])
		}
	})
}
