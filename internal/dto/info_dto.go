package dto

type InventoryItemResp struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ReceivedTxResp struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTxResp struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type CoinHistoryResp struct {
	Received []ReceivedTxResp `json:"received"`
	Sent     []SentTxResp     `json:"sent"`
}

type InfoResponse struct {
	Coins       int                 `json:"coins"`
	Inventory   []InventoryItemResp `json:"inventory"`
	CoinHistory CoinHistoryResp     `json:"coinHistory"`
}
