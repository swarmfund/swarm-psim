package ethsupervisor

import (
	"math/big"

	"gitlab.com/swarmfund/psim/psim/supervisor"
)

type Config struct {
	Supervisor    supervisor.Config `fig:"supervisor"`
	Confirmations uint64
	Cursor        *big.Int
	BaseAsset     string
	DepositAsset  string
}
