package utils

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/tokend/keypair"
)

var (
	ETHHooks = figure.Hooks{
		// TODO move upstream
		"uint64": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToUint64E(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse uint64")
			}
			return reflect.ValueOf(result), nil
		},
		"common.Address": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				if !common.IsHexAddress(v) {
					// provide value does not look like valid address
					return reflect.Value{}, errors.New("invalid address")
				}
				return reflect.ValueOf(common.HexToAddress(v)), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
		"*common.Address": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				if !common.IsHexAddress(v) {
					// provide value does not look like valid address
					return reflect.Value{}, errors.New("invalid address")
				}
				addr := common.HexToAddress(v)
				return reflect.ValueOf(&addr), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
		"keypair.Address": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				kp, err := keypair.ParseAddress(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse kp")
				}
				return reflect.ValueOf(kp), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
		"keypair.Full": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				kp, err := keypair.ParseSeed(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse kp")
				}
				kpFull, ok := kp.(keypair.Full)
				if !ok {
					return reflect.Value{}, errors.Wrap(err,
						"failed to cast kp to keypair.Full; string must be a Seed")
				}
				return reflect.ValueOf(kpFull), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)
