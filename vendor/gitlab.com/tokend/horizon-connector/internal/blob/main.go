package blob

import (
	"gitlab.com/tokend/horizon-connector/internal"
	"fmt"
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector/internal/resources"
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
	url := fmt.Sprintf("/blobs/%s", blobID)
	fields := logan.F{
		"request_url": url,
	}

	respBB, err := q.client.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send GET request", fields)
	}

	if respBB == nil {
		// No such Blob
		return nil, nil
	}
	fields["raw_response"] = string(respBB)

	var response struct {
		Data resources.Blob `json:"data"`
	}
	if err := json.Unmarshal(respBB, &response); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response bytes", fields)
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
	if err != nil {
		return "", errors.Wrap(err, "Failed to send request")
	}
	fields := logan.F{
		"raw_response": string(respBB),
	}

	var respBlob struct{
		Data resources.Blob `json:"data"`
	}
	err = json.Unmarshal(respBB, &respBlob)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal response bytes into Blob struct", fields)
	}

	return respBlob.Data.ID, nil
}
