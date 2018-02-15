package account

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/internal"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/responses"
)

type Q struct {
	client internal.Client
}

func NewQ(client internal.Client) *Q {
	return &Q{
		client,
	}
}

func (q *Q) Signers(address string) ([]resources.Signer, error) {
	endpoint := fmt.Sprintf("/accounts/%s/signers", address)
	response, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}
	var result responses.AccountSigners
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return result.Signers, nil
}

// ByBalance return account address by balance
// TODO probably move to balances q
func (q *Q) ByBalance(balanceID string) (*string, error) {
	endpoint := fmt.Sprintf("/balances/%s/account", balanceID)
	response, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	// actually it's different struct (HistoryAccount) but it works since we only need account_id
	var account resources.Account
	if err := json.Unmarshal(response, &account); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &account.AccountID, nil
}

func (q *Q) ByAddress(address string) (*resources.Account, error) {
	endpoint := fmt.Sprintf("/accounts/%s", address)
	response, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var account resources.Account
	if err := json.Unmarshal(response, &account); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &account, nil
}

func (q *Q) Balances(address string) ([]resources.Balance, error) {
	endpoint := fmt.Sprintf("/accounts/%s/balances", address)
	response, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var result responses.AccountBalances
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return result, nil
}
