package utils

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/psim/figure"
)

var (
	CommonHooks = figure.Hooks{
		"keypair.KP": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				kp, err := keypair.Parse(v)
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

		"*keypair.Full": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				kp, err := keypair.Parse(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse kp")
				}
				kpFull, ok := kp.(*keypair.Full)
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
		"time.Duration": func(value interface{}) (reflect.Value, error) {
			str, err := cast.ToStringE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse string")
			}
			result, err := time.ParseDuration(str)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse duration")
			}
			return reflect.ValueOf(result), nil
		},
		// ToDo: Move to psim/figure/BaseHooks
		"*url.URL": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				u, err := url.Parse(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse url")
				}
				return reflect.ValueOf(u), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)
