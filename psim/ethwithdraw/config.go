package ethwithdraw

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	// hex encoded private key
	Key      string
	Asset    string
	GasPrice *big.Int
	Source   keypair.Address
	Signer   keypair.Full
	// Token contract address
	Token *common.Address
}
