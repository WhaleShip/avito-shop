package dto

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int64  `json:"amount"`
}
