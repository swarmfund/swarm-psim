package postage

import (
	"fmt"

	"github.com/spf13/viper"
)

type Conf struct {
	Key  string
	From string
}

var conf *Conf

func (c Conf) Namespace() string {
	return "postage"
}

func (c Conf) Init(v *viper.Viper) error {
	fromName := v.GetString("from_name")
	fromEmail := v.GetString("from_email")

	conf = &Conf{
		Key:  v.GetString("key"),
		From: fmt.Sprintf("%s <%s>", fromName, fromEmail),
	}
	return nil
}
