package internal

import (
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	types.Transaction
	BlockNumber uint64
	Timestamp   time.Time
}
