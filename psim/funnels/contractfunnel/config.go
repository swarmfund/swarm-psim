package contractfunnel

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/utils"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	GasPrice        *big.Int         `fig:"gas_price,required"`
	Threshold       *big.Int         `fig:"threshold,required"`
	HotWallet       common.Address   `fig:"hot_wallet,required"`
	PrivateKey      string           `fig:"private_key,required"`
	Tokens          []common.Address `fig:"tokens,required"`
	ExternalSystems []int32          `fig:"external_systems,required"`
	Signer          keypair.Full     `fig:"signer,required"`
}

func NewConfig(raw map[string]interface{}) (config Config, err error) {
	err = figure.Out(&config).From(raw).With(utils.ETHHooks, hooks, figure.BaseHooks).Please()
	if err != nil {
		return config, errors.Wrap(err, "failed to figure out config")
	}
	return config, err
}

var hooks = figure.Hooks{
	"[]int32": func(raw interface{}) (reflect.Value, error) {
		var result []int32
		slice, err := cast.ToIntSliceE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "failed to to cast []int")
		}
		for _, i := range slice {
			result = append(result, int32(i))
		}
		return reflect.ValueOf(result), nil
	},
	"[]common.Address": func(raw interface{}) (reflect.Value, error) {
		addressesStrings, err := cast.ToStringSliceE(raw)
		if err != nil {
			return reflect.Value{}, errors.Wrap(err, "Failed to cast provider to map[string]interface{}")
		}

		var addresses []common.Address
		for i, addrStr := range addressesStrings {
			if !common.IsHexAddress(addrStr) {
				// provide value does not look like valid address
				return reflect.Value{}, errors.From(errors.New("invalid address"), logan.F{
					"address_string": addrStr,
					"address_i":      i,
				})
			}

			addresses = append(addresses, common.HexToAddress(addrStr))
		}

		return reflect.ValueOf(addresses), nil
	},
}
