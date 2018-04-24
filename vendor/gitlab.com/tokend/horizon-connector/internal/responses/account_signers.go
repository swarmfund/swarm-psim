package responses

import "gitlab.com/tokend/horizon-connector/internal/resources"

type AccountSigners struct {
	Signers []resources.Signer `json:"signers"`
}
