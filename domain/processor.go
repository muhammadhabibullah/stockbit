package domain

import (
	"context"
)

type Processor interface {
	Run(ctx context.Context)
}
