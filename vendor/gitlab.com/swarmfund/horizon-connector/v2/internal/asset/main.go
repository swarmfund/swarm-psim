package asset

import (
	"fmt"

	"net/http"

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
	resp, err := q.client.Get(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	case http.StatusOK:
		var asset resources.Asset
		if err := json.NewDecoder(resp.Body).Decode(&asset); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal")
		}
		return &asset, nil
	default:
		return nil, errors.Wrapf(err, "request failed with %d", resp.StatusCode)
	}
}
