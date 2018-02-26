package btcdepositveri

import "gitlab.com/tokend/keypair"

type Config struct {
	Host string `fig:"host"`
	Port int    `fig:"port"`

	DepositAsset       string `fig:"deposit_asset"`
	MinDepositAmount   uint64 `fig:"min_deposit_amount"`
	FixedDepositFee    uint64 `fig:"fixed_deposit_fee"`
	LastBlocksNotWatch uint64 `fig:"last_blocks_not_watch"`

	Signer keypair.Full `fig:"signer" mapstructure:"signer"`
}
