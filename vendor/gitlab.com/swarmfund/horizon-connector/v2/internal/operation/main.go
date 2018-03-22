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

func (q *Q) AllRequests(cursor string) ([]resources.Request, error) {
	url := fmt.Sprintf("/requests?limit=200&cursor=%s", cursor)
	return q.getRequests(url)
}

// DEPRECATED Instead use Requests providing specific ReviewableRequestType
func (q *Q) WithdrawalRequests(cursor string) ([]resources.Request, error) {
	url := fmt.Sprintf("/request/withdrawals?limit=200&cursor=%s", cursor)
	return q.getRequests(url)
}

// Requests obtains batch of Requests of the provided type from the provided cursor
// It differs from the AllRequests method, as the latter uses /requests path to obtain Requests.
func (q *Q) Requests(cursor string, reqType ReviewableRequestType) ([]resources.Request, error) {
	url := fmt.Sprintf("/request/%s?limit=200&cursor=%s", string(reqType), cursor)
	return q.getRequests(url)
}

func (q *Q) getRequests(url string) ([]resources.Request, error) {
	response, err := q.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed", logan.F{
			"request_url": url,
		})
	}

	var result responses.RequestsIndex
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response", logan.F{
			"raw_response": string(response),
			"request_url": url,
		})
	}

	return result.Embedded.Records, nil
}
