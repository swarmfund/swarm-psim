package btcdeposit

import (
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Supervisor supervisor.Config `fig:"supervisor"`

	LastProcessedBlock uint64 `fig:"last_processed_block"`
	LastBlocksNotWatch uint64 `fig:"last_blocks_not_watch"`
	MinDepositAmount   uint64 `fig:"min_deposit_amount"`
	DepositAsset       string `fig:"deposit_asset"`
	FixedDepositFee    uint64 `fig:"fixed_deposit_fee"`

	Signer keypair.Full    `fig:"signer" mapstructure:"signer"`
	Source keypair.Address `fig:"exchange"`
}
