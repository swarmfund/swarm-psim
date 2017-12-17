package resources

import "gitlab.com/swarmfund/go/xdr"

type Signer struct {
	PublicKey string         `json:"public_key"`
	Weight    int32          `json:"weight"`
	Type      xdr.SignerType `json:"signer_type"`
	Identity  int32          `json:"signer_identity"`
	Name      string         `json:"signer_name"`
}
