package taxman

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cast"

	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/psim/figure"
)

type Config struct {
	Signer        keypair.KP
	ServiceName   string
	Source        keypair.KP
	Pprof         bool
	Host          string
	Port          int
	LeadershipKey string
	Skip          SkipTransactions
	// DisableVerify will force leader to not attempt to sync snapshoter,
	// no transactions will be submitted
	DisableVerify bool
}

type SkipTransactions map[string]bool

var (
	SkipTransactionsHook = figure.Hooks{
		"taxman.SkipTransactions": func(value interface{}) (reflect.Value, error) {
			result := SkipTransactions{}
			slice, err := cast.ToStringSliceE(value)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to parse []string")
			}
			for _, tx := range slice {
				result[tx] = true
			}
			return reflect.ValueOf(result), nil
		},
	}
)
