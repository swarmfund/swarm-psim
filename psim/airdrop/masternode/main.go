package masternode

import (
	"context"
	"math/big"
	"time"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/funnels/contractfunnel"
	"gitlab.com/swarmfund/psim/psim/internal/eth"
	"gitlab.com/swarmfund/psim/psim/utils"
)

type Config struct {
	FirstBlockTime time.Time      `fig:"first_block_time,required"`
	BlockDuration  time.Duration  `fig:"block_duration,required"`
	Issuer         eth.Keypair    `fig:"issuer,required"`
	Token          common.Address `fig:"token,required"`
	PromoteAfter   int64          `fig:"promote_after,required"`
	Reward         *big.Int       `fig:"reward,required"`
	Furnace        common.Address `fig:"furnace,required"`
	GasPrice       *big.Int       `fig:"gas_price,required"`
	GasLimit       uint64         `fig:"gas_limit,required"`
}

func NewConfig(raw map[string]interface{}) (*Config, error) {
	var config Config
	err := figure.
		Out(&config).
		From(raw).
		With(figure.BaseHooks, eth.KeypairHook, utils.ETHHooks).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}
	return &config, nil
}

func init() {
	app.RegisterService(conf.ServiceAirdropMasternode, func(ctx context.Context) (app.Service, error) {
		geth := app.Config(ctx).Ethereum()

		config, err := NewConfig(app.Config(ctx).Get(conf.ServiceAirdropMasternode))
		if err != nil {
			return nil, errors.Wrap(err, "failed to init config")
		}

		// issuer nonce expected to equal paid out blocks
		nonce, err := geth.PendingNonceAt(ctx, config.Issuer.Address())
		if err != nil {
			return nil, errors.Wrap(err, "failed to get issuer nonce")
		}

		ns := NodesState{
			eth:   geth,
			token: config.Token,
			rq: RewardQueue{
				promoteAfter: config.PromoteAfter,
				blacklist:    []string{config.Issuer.Address().Hex()},
			},
			furnace:       config.Furnace,
			firstBlock:    config.FirstBlockTime,
			blockDuration: config.BlockDuration,
			currentBlock:  nonce,
		}

		contract, err := contractfunnel.NewERC20(config.Token, geth)
		if err != nil {
			panic(err)
		}

		for {
			destination, meta, err := ns.Payout(ctx)
			if err != nil {
				panic(err)
			}
			app.Log(ctx).WithFields(logan.F{
				"block":       meta.Block,
				"destination": destination.Hex(),
				"block_time":  meta.BlockTime.String(),
			}).Info("payout ready")
			_, err = contract.Transfer(&bind.TransactOpts{
				From:  config.Issuer.Address(),
				Nonce: big.NewInt(0).SetUint64(meta.Block),
				Signer: func(_ types.Signer, _ common.Address, tx *types.Transaction) (*types.Transaction, error) {
					return config.Issuer.SignTX(tx)
				},
				GasPrice: eth.FromGwei(config.GasPrice),
				GasLimit: config.GasLimit,
			}, destination, eth.FromGwei(config.Reward))
			if err != nil {
				panic(err)
			}
		}
	})
}
