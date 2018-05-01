package metrics

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Config struct {
	Host string `fig:"host"`
	Port int    `fig:"port,required"`
}

func NewConfig(configData map[string]interface{}) (*Config, error) {
	config := &Config{
		Host: "localhost",
	}

	err := figure.
		Out(config).
		From(configData).
		Please()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to figure out")
	}

	return config, nil
}
