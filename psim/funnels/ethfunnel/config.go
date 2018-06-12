package ethfunnel

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Seed          string         `fig:"seed,required"`
	Destination   common.Address `fig:"destination,required"`
	Threshold     *big.Int       `fig:"threshold,required"`
	GasPrice      *big.Int       `fig:"gas_price,required"`
	Confirmations *big.Int       `fig:"confirmations"`
	// AccountsToDerive will tell wallet how many keys it should derive
	AccountsToDerive uint64 `fig:"accounts_to_derive,required"`
}
