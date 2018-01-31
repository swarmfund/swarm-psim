package btcsupervisor

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

	// TODO Remove after moving logic of the second signature to the verifier.
	AdditionalSignerKP keypair.Full `fig:"additional_signer" mapstructure:"signer"`
}
