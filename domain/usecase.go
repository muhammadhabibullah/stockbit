package domain

import (
	"context"
)

type UserUseCase interface {
	Deposit(ctx context.Context, req DepositRequest) error
	GetDeposit(ctx context.Context, walletID int64) (*GetDepositResponse, error)
}
