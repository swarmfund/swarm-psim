package conf

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/figure"
)

func (c *ViperConfig) Ethereum() *ethclient.Client {
	config := struct {
		Host string
		Port int
	}{
		Host: "localhost",
		Port: 8545,
	}
	err := figure.
		Out(&config).
		From(c.Get("ethereum")).
		With(figure.BaseHooks).
		Please()
	if err != nil {
		panic("failed to figure out ethereum")
	}

	client, err := ethclient.Dial(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
	if err != nil {
		panic(errors.Wrap(err, "failed to dial eth"))
	}
	return client
}
