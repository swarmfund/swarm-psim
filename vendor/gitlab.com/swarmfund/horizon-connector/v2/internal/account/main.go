package account

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	resp, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var result responses.AccountSigners
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		return result.Signers, nil
	default:
		return nil, errors.Wrapf(err, "request failed with %d", resp.StatusCode)
	}
}

func (q *Q) ByAddress(address string) (*resources.Account, error) {
	endpoint := fmt.Sprintf("/accounts/%s", address)
	resp, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var account resources.Account
		if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		return &account, nil
	default:
		return nil, errors.Wrapf(err, "request failed with %d", resp.StatusCode)
	}
}

func (q *Q) Balances(address string) ([]resources.Balance, error) {
	endpoint := fmt.Sprintf("/accounts/%s/balances", address)
	resp, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var result responses.AccountBalances
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		return result, nil
	default:
		return nil, errors.Wrapf(err, "request failed with %d", resp.StatusCode)
	}
}
