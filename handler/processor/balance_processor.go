package processor

import (
	"log"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"
	"stockbit/config"
	"stockbit/domain/proto/pb"
)

type balanceProcessor struct {
	debug bool
}

func NewBalanceProcessor(cfg config.Config) *balanceProcessor {
	return &balanceProcessor{
		debug: cfg.Server.Debug,
	}
}

func (p *balanceProcessor) Handle(ctx goka.Context, msg any) {
	balance := new(pb.Balance)
	if val := ctx.Value(); val != nil {
		balance = val.(*pb.Balance)
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

	balance.Amount += deposit.Amount

	ctx.SetValue(balance)
	if p.debug {
		log.Printf("topic %s; %+v", ctx.Key(), balance)
	}
}
