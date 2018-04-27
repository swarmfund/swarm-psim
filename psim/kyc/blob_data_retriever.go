package kyc

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
)

const (
	// const starts with KYC intentionally
	KYCFormBlobType = "kyc_form"
)

type BlobsConnector interface {
	Blob(blobID string) (*horizon.Blob, error)
}

type BlobDataRetriever struct {
	blobsConnector BlobsConnector
}

func NewBlobDataRetriever(connector BlobsConnector) *BlobDataRetriever {
	return &BlobDataRetriever{}
}

// ParseBlobData retrieves KYC Blob and parses KYCData from Blob.Attributes.Value.
func (p *BlobDataRetriever) ParseBlobData(kycRequest horizon.KYCRequest, accountID string) (*Data, error) {
	blob, err := p.RetrieveKYCBlob(kycRequest, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve KYC Blob")
	}

	kycData, err := ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse KYCData form the Blob.Attributes.Value")
	}

	return kycData, nil
}

// RetrieveKYCBlob retrieves BlobID from KYCRequest,
// obtains Blob by BlobID,
// check Blob's type
// and returns Blob if everything's fine.
func (p *BlobDataRetriever) RetrieveKYCBlob(kycRequest horizon.KYCRequest, accountID string) (*horizon.Blob, error) {
	blobIDInterface, ok := kycRequest.KYCData["blob_id"]
	if !ok {
		return nil, errors.New("Cannot found 'blob_id' key in the KYCData map in the KYCRequest.")
	}
	fields := logan.F{
		"blob_id": blobIDInterface,
	}

	blobID, ok := blobIDInterface.(string)
	if !ok {
		// Normally should never happen
		return nil, errors.From(errors.New("BlobID from KYCData map of the KYCRequest is not a string."), fields)
	}
	fields["blob_id"] = blobID

	blob, err := p.blobsConnector.Blob(blobID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Blob from Horizon", fields)
	}
	if blob == nil {
		return nil, errors.From(errors.New("Could not find Blob by BlobID in Horizon."), fields)
	}
	fields["blob"] = blob

	if blob.Type != KYCFormBlobType {
		return nil, errors.From(errors.Errorf("The Blob provided in KYC Request is of type (%s), but expected (%s).",
			blob.Type, KYCFormBlobType), fields)
	}

	return blob, nil
}
