package domain

type DepositRequest struct {
	WalletID int64   `json:"wallet_id"`
	Amount   float64 `json:"amount"`
}

type GetDepositResponse struct {
	Amount         float64 `json:"amount"`
	AboveThreshold bool    `json:"above_threshold"`
}
