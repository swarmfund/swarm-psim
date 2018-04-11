package ethcontracts

import (
	"math/big"

	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source keypair.Address `fig:"source,required"`
	Signer keypair.Full    `fig:"signer,required" mapstructure:"signer"`

	ETHPrivateKey string   `fig:"eth_private_key,required"`
	GasPrice      *big.Int `fig:"gas_price,required"`
}
