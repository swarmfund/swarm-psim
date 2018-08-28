package eth_deposit

import (
	"time"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	Source          keypair.Address `fig:"source,required"`
	Signer          keypair.Full    `fig:"signer,required"`
	Keypair         eth.Keypair     `fig:"keypair,required"`
	Asset           string          `fig:"asset,required"`
	DisableDeposit  bool            `fig:"disable_deposit"`
	DisableWithdraw bool            `fig:"disable_withdraw"`
	PollingTimeout  time.Duration   `fig:"polling_timeout"`
}

func NewConfig(raw map[string]interface{}) (config Config, err error) {
	err = figure.
		Out(&config).
		From(raw).
		With(figure.BaseHooks, eth.KeypairHook, utils.ETHHooks).
		Please()
	return config, err
}
