package responses

import "gitlab.com/swarmfund/horizon-connector/v2/internal/resources"

type AccountSigners struct {
	Signers []resources.Signer `json:"signers"`
}
