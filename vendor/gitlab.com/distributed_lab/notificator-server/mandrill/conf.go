package mandrill

import "github.com/spf13/viper"

type Conf struct {
	API       string
	Key       string
	FromEmail string
	FromName  string
}

var conf *Conf

func GetConf() *Conf {
	return conf
}

func (c Conf) Namespace() string {
	return "mandrill"
}

func (c Conf) Init(v *viper.Viper) error {
	conf = &Conf{
		API: v.GetString("api"),
		Key: v.GetString("key"),
		FromEmail: v.GetString("from_email"),
		FromName: v.GetString("from_name"),
	}
	return nil
}