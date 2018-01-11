package ethwithdraw

import (
	"math/big"

	"gitlab.com/swarmfund/go/keypair"
)

type Config struct {
	// hex encoded private key
	Key      string
	Asset    string
	GasPrice *big.Int
	Source   keypair.KP
	Signer   keypair.KP
}
