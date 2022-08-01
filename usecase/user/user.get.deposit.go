package user

import (
	"context"
	"fmt"
	"time"

	"stockbit/domain"
	"stockbit/domain/proto/pb"
)

func (u *userUseCase) GetDeposit(_ context.Context, walletID int64) (*domain.GetDepositResponse, error) {
	key := fmt.Sprintf("%d", walletID)

	balanceView, err := u.viewers[domain.BalanceGroup].Get(key)
	if err != nil {
		return nil, err
	}
	if balanceView == nil {
		return &domain.GetDepositResponse{}, nil
	}

	balance, ok := balanceView.(*domain.Balance)
	if !ok {
		return nil, fmt.Errorf("unsupported balance type: %T", balanceView)
	}

	aboveThresholdView, err := u.viewers[domain.AboveThresholdGroup].Get(key)
	if err != nil {
		return nil, err
	}

	aboveThreshold, ok := aboveThresholdView.(*pb.AboveThresholdTable)
	if !ok {
		return nil, fmt.Errorf("unsupported aboveThreshold type: %T", aboveThresholdView)
	}

	var (
		totalBalance         float64
		thresholdCfg         = u.depositCfg.Threshold
		totalBalanceLimit    = thresholdCfg.Amount
		timeLimitDuration, _ = time.ParseDuration(thresholdCfg.Time)
		timeLimit            = time.Now().Add(-timeLimitDuration)
	)

	for _, balanceHistory := range aboveThreshold.BalanceHistory {
		if balanceHistory.CreatedAt.AsTime().Before(timeLimit) {
			break
		}

		totalBalance += float64(balanceHistory.Amount)
		if totalBalance > totalBalanceLimit {
			break
		}
	}

	return &domain.GetDepositResponse{
		Amount:         balance.Amount,
		AboveThreshold: totalBalance > totalBalanceLimit,
	}, nil
}
