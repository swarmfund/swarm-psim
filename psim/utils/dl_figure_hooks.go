package utils

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/distributed_lab/figure"
)

var (
	ETHHooks = figure.Hooks{
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
	}
)
