package airdrop

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
)

var (
	ErrNoBalanceID  = errors.New("BalanceID not found for Account.")
)

type AccountsConnector interface {
	Balances(address string) ([]horizon.Balance, error)
}

type BalanceIDProvider struct {
	accConnector AccountsConnector
}

func NewBalanceIDProvider(accConnector AccountsConnector) *BalanceIDProvider {
	return &BalanceIDProvider{
		accConnector: accConnector,
	}
}

// GetBalanceID retrieves all Balances of the Account and finds the Balance
// of the provided asset among them.
//
// If Account does not have a Balance for the provided asset - (nil, nil) will be returned.
func (p *BalanceIDProvider) GetBalanceID(accAddress, asset string) (*string, error) {
	balances, err := p.accConnector.Balances(accAddress)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Account Balances")
	}

	for _, b := range balances {
		if b.Asset == asset {
			return &b.BalanceID, nil
		}
	}

	return nil, nil
}

// DEPRECATED use BalanceIDProvider instead
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
