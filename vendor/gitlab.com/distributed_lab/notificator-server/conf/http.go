package conf

import "github.com/spf13/viper"

type HTTPConf struct {
	Host string
	Port int
	AllowUntrusted bool
}

func (c HTTPConf) Namespace() string {
	return "http"
}

func (c HTTPConf) Init(v *viper.Viper) error {
	conf := HTTPConf{
		Host: viper.GetString("http.host"),
		Port: viper.GetInt("http.port"),
		AllowUntrusted: viper.GetBool("http.allow_untrusted"),
	}

	httpConf = &conf

	return nil
}

var httpConf *HTTPConf

func GetHTTPConf() *HTTPConf {
	return httpConf
}