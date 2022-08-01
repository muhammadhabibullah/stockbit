package user

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	"stockbit/domain"
	"stockbit/domain/proto/pb"
)

func (u *userUseCase) Deposit(_ context.Context, req domain.DepositRequest) error {
	depositEmit := pb.Deposit{
		WalletId: req.WalletID,
		Amount:   float32(req.Amount),
	}

	msg, err := proto.Marshal(&depositEmit)
	if err != nil {
		return err
	}

	key := fmt.Sprint(req.WalletID)
	value := string(msg)

	err = u.emitters[domain.DepositsTopic].EmitSync(key, value)
	if err != nil {
		return err
	}

	return nil
}
