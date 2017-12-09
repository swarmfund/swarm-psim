package twilio

import "github.com/spf13/viper"

type Conf struct {
	SID        string
	Token      string
	FromNumber string
}

func (c Conf) Namespace() string {
	return "twilio"
}

func (c Conf) Init(v *viper.Viper) error {
	conf = &Conf{
		SID: v.GetString("sid"),
		Token: v.GetString("token"),
		FromNumber: v.GetString("from_number"),
	}
	return nil
}

var conf *Conf