package operation

import (
	"fmt"

	"encoding/json"

	"gitlab.com/swarmfund/horizon-connector/v2/internal"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/responses"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Q struct {
	client internal.Client
}

func NewQ(client internal.Client) *Q {
	return &Q{
		client,
	}
}

// DEPRECATED Does not work anymore.
func (q *Q) Requests(cursor string) ([]resources.Request, error) {
	response, err := q.client.Get(fmt.Sprintf("/requests?cursor=%s", cursor))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	var result responses.RequestsIndex
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return result.Embedded.Records, nil
}

func (q *Q) WithdrawalRequests(cursor string) ([]resources.Request, error) {
	url := fmt.Sprintf("/request/withdrawals?cursor=%s", cursor)

	response, err := q.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed", logan.F{"request_url": url})
	}

	var result responses.RequestsIndex
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response", logan.F{"response": string(response)})
	}

	return result.Embedded.Records, nil
}
