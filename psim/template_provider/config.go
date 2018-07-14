package template_provider

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

type Config struct {
	Host               string `fig:"host"`
	Port               int    `fig:"port"`
	Bucket             string `fig:"bucket,required"`
	SkipSignatureCheck bool   `fig:"skip_signature_check"`
}

func NewConfig(raw map[string]interface{}) (*Config, error) {
	config := Config{
		Host: "localhost",
		Port: 2323,
	}
	err := figure.Out(&config).From(raw).Please()
	if err != nil {
		return nil, errors.Wrap(err, "failed to figure out")
	}
	return &config, nil
}
