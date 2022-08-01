package processor

import (
	"log"
	"time"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"stockbit/config"
	"stockbit/domain/proto/pb"
)

type aboveThresholdProcessor struct {
	debug bool
}

func NewAboveThresholdProcessor(cfg config.Config) *aboveThresholdProcessor {
	return &aboveThresholdProcessor{
		debug: cfg.Server.Debug,
	}
}

func (p *aboveThresholdProcessor) Handle(ctx goka.Context, msg any) {
	aboveThreshold := &pb.AboveThresholdTable{
		BalanceHistory: []*pb.Balance{},
	}
	if val := ctx.Value(); val != nil {
		aboveThreshold = val.(*pb.AboveThresholdTable)
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

	createdAt := ctx.Timestamp()
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	aboveThreshold.BalanceHistory = append([]*pb.Balance{
		{
			Amount:    deposit.Amount,
			CreatedAt: timestamppb.New(createdAt),
		},
	}, aboveThreshold.BalanceHistory...)

	ctx.SetValue(aboveThreshold)
	if p.debug {
		log.Printf("topic %s; %+v", ctx.Key(), aboveThreshold)
	}
}
