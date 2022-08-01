package domain

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"stockbit/domain/proto/pb"
)

type AboveThresholdCodec struct{}

func (atc *AboveThresholdCodec) Encode(value interface{}) ([]byte, error) {
	at, isAboveThreshold := value.(*pb.AboveThresholdTable)
	if !isAboveThreshold {
		return nil, fmt.Errorf("codec requires value *AboveThreshold, got %T", value)
	}

	return proto.Marshal(at)
}

func (atc *AboveThresholdCodec) Decode(data []byte) (interface{}, error) {
	var (
		at  pb.AboveThresholdTable
		err error
	)

	err = proto.Unmarshal(data, &at)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling AboveThreshold: %v", err)
	}

	return &at, nil
}
