package documnet

import (
	"gitlab.com/swarmfund/horizon-connector/v2/internal"
	"fmt"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/distributed_lab/logan/v3"
)

type Q struct {
	client internal.Client
}

func NewQ(client internal.Client) *Q {
	return &Q{
		client,
	}
}

// Document obtains a single Document by its ID.
// If Document doesn't exist - nil,nil is returned.
func (q *Q) Document(docID string) (*resources.Document, error) {
	respBB, err := q.client.Get(fmt.Sprintf("/documents/%s", docID))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send GET request")
	}

	if respBB == nil {
		// No such Document
		return nil, nil
	}

	document := resources.Document{}
	if err := json.Unmarshal(respBB, &document); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response bytes", logan.F{
			"raw_response": string(respBB),
		})
	}

	return &document, nil
}
