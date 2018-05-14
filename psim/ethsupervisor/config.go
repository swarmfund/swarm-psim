package ethsupervisor

import (
	"math/big"

	"gitlab.com/swarmfund/psim/psim/supervisor"
)

type Config struct {
	Supervisor      supervisor.Config `fig:"supervisor,required"`
	Confirmations   uint64            `fig:"confirmations"`
	Cursor          *big.Int          `fig:"cursor"`
	BaseAsset       string            `fig:"base_asset,required"`
	DepositAsset    string            `fig:"deposit_asset,required"`
	FixedDepositFee *big.Int          `fig:"fixed_deposit_fee,required"`
}
