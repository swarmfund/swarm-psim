package utils

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/tokend/go/amount"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/regources"
)

var (
	CommonHooks = figure.Hooks{
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

		"regources.Amount": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				int64Value, err := amount.Parse(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse int64 value")
				}

				return reflect.ValueOf(regources.Amount(int64Value)), nil
			case int:
				return reflect.ValueOf(regources.Amount(int64(v))), nil
			case int64:
				return reflect.ValueOf(regources.Amount(v)), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)
