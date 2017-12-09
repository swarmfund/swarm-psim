package internal

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	types.Transaction
	Timestamp time.Time
}
