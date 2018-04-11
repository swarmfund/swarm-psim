package ethfunnel

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Seed        string
	Destination common.Address
	// Min withdraw amount
	Threshold     *big.Int
	GasPrice      *big.Int
	Confirmations *big.Int
	// AccountsToDerive will tell wallet how many keys it should derive
	AccountsToDerive uint64
}
