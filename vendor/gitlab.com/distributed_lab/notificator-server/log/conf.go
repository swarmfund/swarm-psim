package log

import (
	"github.com/spf13/viper"
	"github.com/Sirupsen/logrus"
	"os"
)

type Conf struct {
	Level string
}

var conf *Conf


func (c Conf) Namespace() string {
	return "log"
}

func (c Conf) Init(v *viper.Viper) error {
	conf := &Conf{
		Level: v.GetString("level"),
	}

	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		panic("failed to parse log level")
	}

	DefaultLogger.Level = level

	return nil
}

func init() {
	// default values before conf takes command
	DefaultLogger = logrus.New()
	DefaultLogger.Level = logrus.FatalLevel
	DefaultEntry = DefaultLogger.WithField("pid", os.Getpid())
}