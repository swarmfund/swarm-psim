package mailgun

import (
	"fmt"

	"github.com/spf13/viper"
)

// Conf is a configuration for Mailgun Client.
type Conf struct {
	Key       string
	PublicKey string
	Domain    string
	From      string
}

var conf *Conf

// GetConf returns mailgun.Conf.
func GetConf() *Conf {
	return conf
}

// Namespace returns name of the mailgun.Conf
// object in the config file.
func (c Conf) Namespace() string {
	return "mailgun"
}

// Init initialize new mailgun.Conf from
// viper.Viper configuration registry
func (c Conf) Init(v *viper.Viper) error {
	fromName := v.GetString("from_name")
	fromEmail := v.GetString("from_email")

	conf = &Conf{
		Key:       v.GetString("key"),
		PublicKey: v.GetString("public_key"),
		Domain:    v.GetString("domain"),
		From:      fmt.Sprintf("%s <%s>", fromName, fromEmail),
	}
	return nil
}
