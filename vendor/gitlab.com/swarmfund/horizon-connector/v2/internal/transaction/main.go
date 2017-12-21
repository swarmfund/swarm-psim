package transaction

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

func (q *Q) Transactions(cursor string) ([]resources.Transaction, error) {
	response, err := q.client.Get(fmt.Sprintf("/transactions?cursor=%s", cursor))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer response.Body.Close()

	var result responses.TransactionIndex
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return result.Embedded.Records, nil
}
