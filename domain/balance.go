package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Balance struct {
	Amount    float64    `json:"amount"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type BalanceCodec struct{}

func (bc *BalanceCodec) Encode(value interface{}) ([]byte, error) {
	if _, isBalance := value.(*Balance); !isBalance {
		return nil, fmt.Errorf("codec requires value *Balance, got %T", value)
	}

	return json.Marshal(value)
}

func (bc *BalanceCodec) Decode(data []byte) (interface{}, error) {
	var (
		b   Balance
		err error
	)

	err = json.Unmarshal(data, &b)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling Balance: %v", err)
	}

	return &b, nil
}
