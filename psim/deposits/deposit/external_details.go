package deposit

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/types"
)

// ExternalDetails is a blob to be put into Issuance Operation ExternalDetails as JSON string.
// ExternalDetails provide info not included into the Issuance itself, but necessary for verification the Issuance.
type ExternalDetails struct {
	BlockNumber uint64       `json:"block_number"`
	TXHash      string       `json:"tx_hash"`
	OutIndex    uint         `json:"out_index"`
	Price       types.Amount `json:"price"`
}

// Encode returns ExternalDetails marshaled into JSON.
func (d ExternalDetails) Encode() string {
	bytes, err := json.Marshal(&d)
	if err != nil {
		panic(errors.Wrap(err, "Failed to encode DepositDetails"))
	}

	return string(bytes)
}

// GetLoganFields implements fields.Provider interface from logan.
func (d ExternalDetails) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"block_number": d.BlockNumber,
		"tx_hash":      d.TXHash,
		"price":        d.Price,
	}
}
