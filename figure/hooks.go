package figure

import (
	"reflect"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	// DEPRECATED Move to using gitlab.com/distributed_lab/figure
	BaseHooks Hooks = Hooks{
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
	}
)

// DEPRECATED Move to using gitlab.com/distributed_lab/figure
type Hook func(value interface{}) (reflect.Value, error)

// DEPRECATED Move to using gitlab.com/distributed_lab/figure
type Hooks map[string]Hook

// Merge does not modify any Hooks, only produces new Hooks.
// If duplicated keys - the value from the last Hooks with such key will be taken.
// DEPRECATED Move to using gitlab.com/distributed_lab/figure
func Merge(manyHooks ...Hooks) Hooks {
	if len(manyHooks) == 1 {
		return manyHooks[0]
	}

	merged := Hooks{}

	for _, hooks := range manyHooks {
		for key, hook := range hooks {
			merged[key] = hook
		}
	}

	return merged
}
