package processor

import (
	"log"

	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"
	"stockbit/domain"
	"stockbit/domain/proto/pb"
)

type balanceProcessor struct{}

func NewBalanceProcessor() *balanceProcessor {
	return &balanceProcessor{}
}

func (p *balanceProcessor) Handle(ctx goka.Context, msg any) {
	balance := new(domain.Balance)
	if val := ctx.Value(); val != nil {
		balance = val.(*domain.Balance)
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

	balance.Amount += float64(deposit.Amount)

	ctx.SetValue(balance)
	log.Printf("topic %s; %+v", ctx.Key(), balance)
}
