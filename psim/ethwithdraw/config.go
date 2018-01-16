package ethwithdraw

import (
	"math/big"

	"gitlab.com/tokend/keypair"
)

type Config struct {
	// hex encoded private key
	Key      string
	Asset    string
	GasPrice *big.Int
	Source   keypair.Address
	Signer   keypair.Full
}
