package models

type User struct {
	Username string
	Password string
	Coins    int64
}

type CoinTransaction struct {
	ID       uint
	FromUser string
	ToUser   string
	Amount   int64
}

type InventoryItem struct {
	ID       uint
	UserName string
	ItemName string
	Quantity int
}

type MerchItem struct {
	ID    uint
	Name  string
	Price int64
}
