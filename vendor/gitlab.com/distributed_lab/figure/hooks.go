package figure

import (
	"reflect"

	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/logan/v3"
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
		"int64": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToInt64E(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse int64")
			}
			return reflect.ValueOf(result), nil
		},
		"uint": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToUintE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse uint")
			}
			return reflect.ValueOf(result), nil
		},
		"uint32": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToUint32E(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse uint32")
			}
			return reflect.ValueOf(result), nil
		},
		"uint64": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToUint64E(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse uint64")
			}
			return reflect.ValueOf(result), nil
		},
		"float64": func(value interface{}) (reflect.Value, error) {
			result, err := cast.ToFloat64E(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse float64")
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
		"*time.Duration": func(value interface{}) (reflect.Value, error) {
			if value == nil {
				return reflect.ValueOf(nil), nil
			}
			result, err := cast.ToDurationE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse duration")
			}
			return reflect.ValueOf(&result), nil
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
		"logan.Level": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				lvl, err := logan.ParseLevel(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse log level")
				}
				return reflect.ValueOf(lvl), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)
