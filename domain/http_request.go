package domain

type DepositRequest struct {
	WalletID int64   `json:"wallet_id" example:"1"`
	Amount   float64 `json:"amount" example:"1000"`
}

type GetDepositResponse struct {
	Amount         float64 `json:"amount" example:"3000"`
	AboveThreshold bool    `json:"above_threshold" example:"false"`
}
