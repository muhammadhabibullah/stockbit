package domain

import (
	"encoding/json"
	"fmt"
)

type AboveThreshold struct {
	BalanceHistory []Balance
}

type AboveThresholdCodec struct{}

func (atc *AboveThresholdCodec) Encode(value interface{}) ([]byte, error) {
	if _, isAboveThreshold := value.(*AboveThreshold); !isAboveThreshold {
		return nil, fmt.Errorf("codec requires value *AboveThreshold, got %T", value)
	}

	return json.Marshal(value)
}

func (atc *AboveThresholdCodec) Decode(data []byte) (interface{}, error) {
	var (
		at  AboveThreshold
		err error
	)

	err = json.Unmarshal(data, &at)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling AboveThreshold: %v", err)
	}

	return &at, nil
}
