package sale

import (
	"encoding/json"
	"fmt"

	"gitlab.com/swarmfund/horizon-connector/v2/internal"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
)

type Q struct {
	client internal.Client
}

func NewQ(client internal.Client) *Q {
	return &Q{
		client,
	}
}

func (q *Q) Sales() ([]resources.Sale, error) {
	endpoint := fmt.Sprintf("/core_sales")
	response, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var result []resources.Sale
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return result, nil
}
