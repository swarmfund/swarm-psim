package figure

import (
	"reflect"

	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

var (
	// BaseHooks set of default hooks for common types
	BaseHooks = Hooks{
		"string": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToStringE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse string")
			}
			return reflect.ValueOf(result), nil
		},
		"[]string": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToStringSliceE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse []string")
			}
			return reflect.ValueOf(result), nil
		},
		"int": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToIntE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse int")
			}
			return reflect.ValueOf(result), nil
		},
		"bool": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToBoolE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse bool")
			}
			return reflect.ValueOf(result), nil
		},
		"*time.Time": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToTimeE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse time")
			}
			return reflect.ValueOf(&result), nil
		},
		"time.Duration": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToDurationE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse duration")
			}
			return reflect.ValueOf(result), nil
		},
		"*big.Int": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				i, ok := new(big.Int).SetString(v, 10)
				if !ok {
					return reflect.Value{}, errors.New("failed to parse")
				}
				return reflect.ValueOf(i), nil
			case int:
				return reflect.ValueOf(big.NewInt(int64(v))), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)
