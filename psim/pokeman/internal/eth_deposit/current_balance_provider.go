package eth_deposit

import (
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CurrentBalanceProvider interface {
	CurrentBalance() (horizon.Balance, error)
}

type currentBalanceProvider struct {
	connector *horizon.Connector
	address   string
	asset     string
}

func NewCurrentBalanceProvider(connector *horizon.Connector, address string, asset string) CurrentBalanceProvider {
	return &currentBalanceProvider{connector, address, asset}
}

func (c *currentBalanceProvider) CurrentBalance() (horizon.Balance, error) {
	balance, err := c.connector.Accounts().CurrentBalanceIn(c.address, c.asset)
	if err != nil {
		return horizon.Balance{}, errors.Wrap(err, "failed to get account balance")
	}
	return balance, nil
}