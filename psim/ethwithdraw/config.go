package ethwithdraw

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/swarmfund/go/keypair"
)

type Config struct {
	Keystore   string
	Passphrase string
	Asset      string
	From       common.Address
	GasPrice   *big.Int
	Source     keypair.KP
	Signer     keypair.KP
}
