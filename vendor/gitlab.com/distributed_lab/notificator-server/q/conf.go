package q

import "github.com/spf13/viper"

type Conf struct {
	Driver string
	DSN string
}

func (c Conf) Init(v *viper.Viper) error {
	c.Driver = v.GetString("driver")
	c.DSN = v.GetString("dsn")

	conf = &c

	return nil
}

func (c Conf) Namespace() string {
	return "db"
}

var conf *Conf