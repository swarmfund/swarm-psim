package data

import (
	"gitlab.com/swarmfund/go/resources"
	"gitlab.com/swarmfund/horizon-connector/v2"
)

type AccountQ struct {
	horizon *horizon.Connector
}

func NewAccountQ(horizon *horizon.Connector) *AccountQ {
	return &AccountQ{
		horizon: horizon,
	}
}

func (q *AccountQ) Signers(address string) ([]resources.Signer, error) {

	signers, err := q.horizon.Accounts().Signers(address)
	if err != nil {
		return nil, err
	}

	if signers == nil {
		return nil, nil
	}

	// TODO share resource
	result := make([]resources.Signer, len(signers))
	for _, signer := range signers {
		result = append(result, resources.Signer{
			AccountID:  signer.PublicKey,
			Weight:     int(signer.Weight),
			SignerType: int(signer.Type),
			Identity:   int(signer.Identity),
			Name:       signer.Name,
		})
	}
	return result, nil
}
