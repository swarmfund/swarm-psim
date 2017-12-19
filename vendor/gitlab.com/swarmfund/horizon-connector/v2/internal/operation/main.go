package operation

import (
	"fmt"

	"encoding/json"

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

func (q *Q) Requests(cursor string) ([]resources.Request, error) {
	response, err := q.client.Get(fmt.Sprintf("/requests?cursor=%s", cursor))
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	defer response.Body.Close()

	var result responses.RequestsIndex
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return result.Embedded.Records, nil
}
