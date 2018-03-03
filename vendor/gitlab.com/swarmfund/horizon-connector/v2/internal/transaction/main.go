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

func (q *Q) Transactions(cursor string) ([]resources.Transaction, *resources.PageMeta, error) {
	response, err := q.client.Get(fmt.Sprintf("/transactions?limit=200&cursor=%s", cursor))
	if err != nil {
		return nil, nil, errors.Wrap(err, "request failed")
	}

	var result responses.TransactionIndex
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal")
	}
	return result.Embedded.Records, &result.Embedded.Meta, nil
}
