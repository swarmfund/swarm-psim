package airdrop

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/horizon-connector"
)

const (
	KYCFormBlobType = "kyc_form"
)

type BlobsConnector interface {
	Blob(blobID string) (*horizon.Blob, error)
}

type USAChecker struct {
	blobsConnector BlobsConnector
}

func NewUSAChecker(blobsConnector BlobsConnector) *USAChecker {
	return &USAChecker{
		blobsConnector: blobsConnector,
	}
}

// TODO Comment
// CheckIsUSA takes BlobID from Account.KYC.Data,
// retrieves Blob from Horizon by BlobID (using BlobsConnector),
// parses KYCData from the retrieved Blob
// and makes USA/notUSA decision based on parsed KYCData.
//
// If Account.KYC.Data is nil - error is returned (it's expected User to be already verified).
func (c *USAChecker) CheckIsUSA(acc horizon.Account) (bool, error) {
	if acc.KYC.Data == nil {
		return false, errors.New("KYCData is nil - could not find KYCBlobID.")
	}
	fields := logan.F{
		"blob_id": acc.KYC.Data.BlobID,
	}

	blob, err := c.blobsConnector.Blob(acc.KYC.Data.BlobID)
	if err != nil {
		return false, errors.Wrap(err, "Failed to get Blob by BlobID", fields)
	}
	fields["blob"] = blob

	if blob.Type != KYCFormBlobType {
		return false, errors.From(errors.Errorf("The Blob provided in KYCData of Account is of type (%s), but expected (%s).",
			blob.Type, KYCFormBlobType), fields)
	}

	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return false, errors.Wrap(err, "Failed tot parse KYC data from Attributes.Value of the Blob", fields)
	}
	fields["kyc_data"] = kycData

	return kycData.IsUSA(), nil
}
