package conf

import (
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config interface {
	Init() error
	HTTP() HTTPConf
	DB() DBConf
	Log() *logrus.Logger
	Mailgun() MailgunConf
	Postage() PostageConf
	Twilio() TwilioConf
	Mandrill() MandrillConf
	Requests() RequestsConf
}

type ViperConfig struct {
	*viper.Viper
	*sync.RWMutex

	// runtime-initialized instances
	db       *DBConf
	http     *HTTPConf
	log      *logrus.Logger
	mandrill *MandrillConf
	mailgun  *MailgunConf
	postage  *PostageConf
	twilio   *TwilioConf
	requests *RequestsConf
}

func NewViperConfig(fn string) Config {
	config := ViperConfig{
		Viper:   viper.GetViper(),
		RWMutex: &sync.RWMutex{},
	}
	config.SetConfigFile(fn)
	return &config
}

func (c *ViperConfig) Init() error {
	c.Lock()
	defer c.Unlock()

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}
	return nil
}
