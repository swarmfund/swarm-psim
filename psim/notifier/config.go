package notifier

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/psim/figure"
	"gitlab.com/tokend/psim/psim/notifier/internal/types"
)

type Config struct {
	ClientUrl   string
	ProjectName string
	Signer      *keypair.Full

	Operations *Operations
	Assets     *Assets

	Pprof bool
	Host  string
	Port  int
}

type Operations struct {
	ClientUrl       string `mapstructure:"client_url"`
	PayloadID       int    `mapstructure:"payload_id"`
	IgnoreOlderThan string `mapstructure:"ignore_older_than"`
	Cursor          string `mapstructure:"cursor"`
}

type Assets struct {
	PayloadID            int          `mapstructure:"payload_id"`
	EmissionThreshold    types.Amount `mapstructure:"emission_threshold"`
	NotificationReceiver string       `mapstructure:"notification_receiver"`
	CheckPeriod          string       `mapstructure:"check_period"`
	Codes                []string     `mapstructure:"codes"`
}

var (
	CommonHooks = figure.Hooks{
		"*notifier.Operations": func(raw interface{}) (reflect.Value, error) {
			result := &Operations{}
			err := mapstructure.Decode(raw, result)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(result), nil
		},
		"*notifier.Assets": func(raw interface{}) (reflect.Value, error) {
			result := &Assets{}
			err := mapstructure.Decode(raw, result)
			if err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(result), nil
		},
	}
)
