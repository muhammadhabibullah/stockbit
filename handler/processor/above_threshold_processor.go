package processor

import (
	"log"
	"time"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"
	"stockbit/domain"
	"stockbit/domain/proto/pb"
)

type aboveThresholdProcessor struct{}

func NewAboveThresholdProcessor() *aboveThresholdProcessor {
	return &aboveThresholdProcessor{}
}

func (p *aboveThresholdProcessor) Handle(ctx goka.Context, msg any) {
	aboveThreshold := &domain.AboveThreshold{
		BalanceHistory: []domain.Balance{},
	}
	if val := ctx.Value(); val != nil {
		aboveThreshold = val.(*domain.AboveThreshold)
	}

	var (
		deposit   pb.Deposit
		msgStr, _ = msg.(string)
		msgBytes  = []byte(msgStr)
	)

	if err := proto.Unmarshal(msgBytes, &deposit); err != nil {
		log.Printf("error unmarshall: %s", err)
		return
	}

	now := time.Now()
	aboveThreshold.BalanceHistory = append([]domain.Balance{
		{
			Amount:    float64(deposit.Amount),
			CreatedAt: &now,
		},
	}, aboveThreshold.BalanceHistory...)

	ctx.SetValue(aboveThreshold)
	log.Printf("topic %s; %+v", ctx.Key(), aboveThreshold)
}
