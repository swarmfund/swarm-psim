package blob

import (
	"gitlab.com/swarmfund/horizon-connector/v2/internal"
	"fmt"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/internal/resources"
	"gitlab.com/distributed_lab/logan/v3"
	"context"
	"bytes"
)

type Q struct {
	client internal.Client
}

func NewQ(client internal.Client) *Q {
	return &Q{
		client,
	}
}

// Blob obtains a single Blob by its ID (hash).
// If Blob doesn't exist - nil,nil is returned.
func (q *Q) Blob(blobID string) (*resources.Blob, error) {
	respBB, err := q.client.Get(fmt.Sprintf("/blobs/%s", blobID))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send GET request")
	}

	if respBB == nil {
		// No such Blob
		return nil, nil
	}

	var response struct {
		Data resources.Blob `json:"data"`
	}
	if err := json.Unmarshal(respBB, &response); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response bytes", logan.F{
			"raw_response": string(respBB),
		})
	}

	return &response.Data, nil
}

func (q *Q) SubmitBlob(ctx context.Context, blobType, attrValue string, relationships map[string]string) (blobID string, err error) {
	blob := resources.Blob {
		Type: blobType,
		Attributes: resources.BlobAttributes{
			Value: attrValue,
		},
	}
	for k, v := range relationships {
		blob.AddRelationship(k, v)
	}

	reqBB, err := json.Marshal(struct{
		Data resources.Blob `json:"data"`
	}{
		Data: blob,
	})
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal request")
	}

	respBB, err := q.client.Post("/blobs", bytes.NewReader(reqBB))
	if err == nil {
		// successful submission
		return "", errors.Wrap(err, "Failed to send request")
	}
	fields := logan.F{
		"raw_response": string(respBB),
	}

	var respBlob resources.Blob
	err = json.Unmarshal(respBB, &respBlob)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal response bytes into Blob struct", fields)
	}

	return respBlob.ID, nil
}
