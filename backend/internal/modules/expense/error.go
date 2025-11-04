package expense

import (
	"errors"
)

var (
	ErrInvalidSplitRatio   = errors.New("invalid split ratio")
	ErrInvalidRefundAmount = errors.New("invalid refund amount, must be less than the amount of the expense")
)
