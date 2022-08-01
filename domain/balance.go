package domain

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"stockbit/domain/proto/pb"
)

type BalanceCodec struct{}

func (bc *BalanceCodec) Encode(value interface{}) ([]byte, error) {
	table, isBalance := value.(*pb.Balance)
	if !isBalance {
		return nil, fmt.Errorf("codec requires value *Balance, got %T", value)
	}

	return proto.Marshal(table)
}

func (bc *BalanceCodec) Decode(data []byte) (interface{}, error) {
	var (
		b   pb.Balance
		err error
	)

	err = proto.Unmarshal(data, &b)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling Balance: %v", err)
	}

	return &b, nil
}
