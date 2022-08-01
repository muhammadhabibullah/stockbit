package domain

type DepositRequest struct {
	WalletID int64   `json:"wallet_id"`
	Amount   float64 `json:"amount"`
}
