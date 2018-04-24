package notifier

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/kyc"
)

type BlobsConnector interface {
	Blob(blobID string) (*horizon.Blob, error)
}

type KYCDataGetter struct {
	blobsConnector BlobsConnector
}

// GetKYCFirstName returns non-nil error, if Blob not found.
func (g *KYCDataGetter) getKYCFirstName(kycData map[string]interface{}) (string, error) {
	blob, err := g.obtainBlob(kycData)
	if err != nil {
		return "", errors.Wrap(err, "Failed to obtain Blob.Attributes.Value by the kycData map")
	}
	fields := logan.F{
		"blob": blob,
	}

	firstName, err := kyc.ParseKYCFirstName(blob.Attributes.Value)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse KYC data from Attributes.Value string in Blob", fields)
	}

	return firstName, nil
}

// ObtainBlob only returns nil-Blob with a non-nil error.
func (g *KYCDataGetter) obtainBlob(kycData map[string]interface{}) (*horizon.Blob, error) {
	blobIDInterface, ok := kycData["blob_id"]
	if !ok {
		return nil, errors.New("'blob_id' key not found in KYCData map")
	}

	blobID, ok := blobIDInterface.(string)
	if !ok {
		return nil, errors.New("BlobID from KYCData map is not a string")
	}
	fields := logan.F{
		"blob_id": blobID,
	}

	blob, err := g.blobsConnector.Blob(blobID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get Blob from Horizon", fields)
	}
	if blob == nil {
		return nil, errors.From(errors.New("Nil Blob received from API"), fields)
	}
	fields["blob"] = blob

	if blob.Type != KYCFormBlobType {
		return nil, errors.From(errors.Errorf("The Blob provided in KYC Request is of type (%s), but expected (%s).",
			blob.Type, KYCFormBlobType), fields)
	}

	return blob, nil
}
