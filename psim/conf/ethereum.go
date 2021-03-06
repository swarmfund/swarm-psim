package conf

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

func (c *ViperConfig) Ethereum() *ethclient.Client {
	config := struct {
		Proto string `fig:"proto"`
		Host  string `fig:"host"`
		Port  int    `fig:"port"`
	}{
		Proto: "http",
		Host:  "localhost",
		Port:  8545,
	}
	err := figure.
		Out(&config).
		From(c.Get("ethereum")).
		With(figure.BaseHooks).
		Please()
	if err != nil {
		panic("failed to figure out ethereum")
	}

	client, err := ethclient.Dial(fmt.Sprintf("%s://%s:%d", config.Proto, config.Host, config.Port))
	if err != nil {
		panic(errors.Wrap(err, "failed to dial eth"))
	}
	return client
}
