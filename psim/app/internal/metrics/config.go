package metrics

type Config struct {
	Host string `fig:"host,required"`
	Port int64  `fig:"port,required"`
}
