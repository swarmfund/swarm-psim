package kyc

import (
	"gitlab.com/tokend/horizon-connector"
)

const (
	// const starts with KYC intentionally
	KYCFormBlobType = "kyc_form"
)

type BlobsConnector interface {
	Blob(blobID string) (*horizon.Blob, error)
}
