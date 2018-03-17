package airdrop

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

var (
	ErrNoBalanceID  = errors.New("BalanceID not found for Account.")
)

// TODO Consider creating BalanceIDProvider, which gets the AccountsConnector in constructor.

type AccountsConnector interface {
	Balances(address string) ([]horizon.Balance, error)
}

func GetBalanceID(accAddress, asset string, accConnector AccountsConnector) (string, error) {
	balances, err := accConnector.Balances(accAddress)
	if err != nil {
		return "", errors.Wrap(err, "Failed to get Account Balances")
	}

	for _, b := range balances {
		if b.Asset == asset {
			return b.BalanceID, nil
		}
	}

	return "", ErrNoBalanceID
}
