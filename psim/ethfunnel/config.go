package ethfunnel

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Seed          string
	Destination   common.Address
	Threshold     *big.Int
	GasPrice      *big.Int
	Confirmations *big.Int
}
