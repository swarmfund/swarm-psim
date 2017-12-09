package server

import (
	"github.com/spf13/viper"
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/log"
	"gitlab.com/distributed_lab/notificator-server/mailgun"
	"gitlab.com/distributed_lab/notificator-server/mandrill"
	"gitlab.com/distributed_lab/notificator-server/postage"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/twilio"
)

type Config interface {
	Namespace() string
	Init(config *viper.Viper) error
}

func InitConf(fn string) {
	entry := log.WithField("service", "conf")
	viper.SetConfigFile(fn)
	err := viper.ReadInConfig()
	if err != nil {
		entry.WithField("reason", err).Fatal("can't read config file")
	}

	configs := []Config{
		log.Conf{},
		q.Conf{},
		conf.HTTPConf{},
		RequestsConf{},
		mandrill.Conf{},
		mailgun.Conf{},
		twilio.Conf{},
		postage.Conf{},
	}

	for _, config := range configs {
		ns := config.Namespace()
		sub := viper.Sub(ns)
		err := config.Init(sub)
		if err != nil {
			entry.WithField("ns", ns).WithField("reason", err).Fatal()
		}
	}
}
