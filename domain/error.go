package domain

import (
	"fmt"
)

var ErrInvalidID = fmt.Errorf("invalid ID")

var ErrInvalidRequestMethod = fmt.Errorf("invalid request method")
