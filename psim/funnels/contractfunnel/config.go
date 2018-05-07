package contractfunnel

import (
	"reflect"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Config struct {
	ETHPrivateKey                string           `fig:"eth_private_key,required"`
	ContractsAddresses           []common.Address `fig:"contracts_addresses,required"`
	TokenReceiverAddress         common.Address   `fig:"tokens_receiver_address,required"`
	TokenToFunnelContractAddress common.Address   `fig:"token_to_funnel_contract_address,required"`
	FunnelPeriod                 time.Duration    `fig:"funnel_period,required"`
	OnlyViewBalances             bool             `fig:"only_view_balances"` // not required, false by default
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"contracts":                        len(c.ContractsAddresses),
		"tokens_receiver_address":          c.TokenReceiverAddress.String(),
		"token_to_funnel_contract_address": c.TokenToFunnelContractAddress.String(),
		"only_view_balances":               c.OnlyViewBalances,
	}
}

var hooks = figure.Hooks{
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
