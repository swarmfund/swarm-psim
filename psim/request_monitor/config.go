package request_monitor

import (
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/keypair"
)

type Config struct {
	//DefaultTimeout present the period after which the request is considered stale
	//Add individual timeouts in Requests.
	DefaultTimeout time.Duration `fig:"default_timeout,required"`
	//SleepPeriod is frequency with which the service checks the requests.
	SleepPeriod time.Duration `fig:"sleep_period"`
	//NotifyPeriod is frequency with which the service notifies certain request.
	NotifyPeriod time.Duration `fig:"notify_period"`
	Signer       keypair.Full  `fig:"signer,required"`
	//String - name of request
	//Requests contain individual timeouts for requests.
	Requests map[string]RequestConfig `fig:"requests"`

	EnableSlack bool `fig:"enable_slack"`
}
type RequestConfig struct {
	Timeout time.Duration
}

var RequestsHook = figure.Hooks{
	"map[string]request_monitor.RequestConfig": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case map[string]interface{}:
			temp := map[string]RequestConfig{}
			conf := RequestConfig{}
			for typeName, val := range v {
				interval, err := cast.ToDurationE(val)
				if err != nil {
					return reflect.Value{}, errors.New("failed to cast")
				}
				conf = temp[typeName]
				conf.Timeout = interval
				temp[typeName] = conf
			}
			return reflect.ValueOf(temp), nil
		default:
			return reflect.Value{}, errors.New(fmt.Sprintf("unsupported conversion from %T", value))
		}
	},
}

func (c Config) Validate() error {
	if c.DefaultTimeout <= 0 {
		return errors.New("default request timeout is invalid, must be bigger then zero")
	}
	for reqType, reqTimeout := range c.Requests {
		if reqTimeout.Timeout <= 0 {
			return errors.Errorf("invalid timeout config for request: %s , must be bigger then zero", reqType)
		}
		ok := false
		for _, rrt := range xdr.ReviewableRequestTypeAll {
			if rrt.ShortString() == reqType {
				ok = true
			}
		}
		if !ok {
			return errors.Errorf("invalid name in request timeout config: %s", reqType)
		}
	}

	return nil
}
