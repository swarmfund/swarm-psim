package notifier

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/notifier/internal/types"
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
	Enable          bool   `mapstructure:"enable"`
	ClientUrl       string `mapstructure:"client_url"`
	PayloadID       int    `mapstructure:"payload_id"`
	IgnoreOlderThan string `mapstructure:"ignore_older_than"`
	Cursor          string `mapstructure:"cursor"`
}

type Assets struct {
	Enable                bool         `mapstructure:"enable"`
	PayloadID             int          `mapstructure:"payload_id"`
	EmissionThreshold     types.Amount `mapstructure:"emission_threshold"`
	CheckPeriod           string       `mapstructure:"check_period"`
	NotificationReceivers []string     `mapstructure:"notification_receivers"`
	Codes                 []string     `mapstructure:"codes"`
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
