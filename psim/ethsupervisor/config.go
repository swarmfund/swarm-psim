package ethsupervisor

import (
	"math/big"

	"gitlab.com/swarmfund/psim/psim/supervisor"
)

type Config struct {
	Supervisor      supervisor.Config `fig:"supervisor"`
	Confirmations   uint64            `fig:"confirmations"`
	Cursor          *big.Int          `fig:"cursor"`
	BaseAsset       string            `fig:"base_asset"`
	DepositAsset    string            `fig:"deposit_asset"`
	FixedDepositFee *big.Int          `fig:"fixed_deposit_fee"`
}
