package deposit

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/types"
)

// TODO Comment
type VerifyRequest struct {
	AccountID string `json:"account_id"`
	Envelope  string `json:"envelope"`
}

// TODO Comment
func (r VerifyRequest) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"account_id": r.AccountID,
		"envelope":   r.Envelope,
	}
}

// TODO Comment
type ExternalDetails struct {
	BlockNumber uint64       `json:"block_number"`
	TXHash      string       `json:"tx_hash"`
	OutIndex    uint         `json:"out_index"`
	Price       types.Amount `json:"price"`
}

// TODO Comment
func (d ExternalDetails) Encode() string {
	bytes, err := json.Marshal(&d)
	if err != nil {
		panic(errors.Wrap(err, "Failed to encode DepositDetails"))
	}
	return string(bytes)
}

// TODO Comment
func (d ExternalDetails) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"block_number": d.BlockNumber,
		"tx_hash":      d.TXHash,
		"price":        d.Price,
	}
}
