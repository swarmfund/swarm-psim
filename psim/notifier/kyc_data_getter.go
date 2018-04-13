package notifier

import (
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type BlobsConnector interface {
	Blob(blobID string) (*horizon.Blob, error)
}

type KYCDataGetter struct {
	blobsConnector BlobsConnector
}

func (g *KYCDataGetter) getBlobKYCData(kycData map[string]interface{}) (*kyc.Data, error) {
	blobIDInterface, ok := kycData["blob_id"]
	if !ok {
		return nil, errors.New("'blob_id' key not found in KYCData map")
	}

	blobID, ok := blobIDInterface.(string)
	if !ok {
		return nil, errors.New("BlobID from KYCData map is not a string")
	}

	blob, err := g.blobsConnector.Blob(blobID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Blob from Horizon")
	}

	fields := logan.F{"blob": blob}

	if blob.Type != KYCFormBlobType {
		return nil, errors.From(errors.Errorf("The Blob provided in KYC Request is of type (%s), but expected (%s).",
			blob.Type, KYCFormBlobType), fields)
	}

	blobKYCData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse KYC data from Attributes.Value string in Blob", fields)
	}

	return blobKYCData, nil
}
