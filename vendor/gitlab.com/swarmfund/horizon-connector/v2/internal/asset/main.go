package asset

import (
	"fmt"

	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/internal"
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
func (q Q) ByCode(code string) (*resources.Asset, error) {
	endpoint := fmt.Sprintf("/assets/%s", code)
	response, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if response == nil {
		return nil, nil
	}

	var asset resources.Asset
	if err := json.Unmarshal(response, &asset); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &asset, nil

}
