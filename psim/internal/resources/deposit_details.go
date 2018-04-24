package resources

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/tokend/horizon-connector/types"
)

// DEPRECATED Use deposit.ExternalDetails instead
type DepositDetails struct {
	TXHash string       `json:"tx_hash"`
	Price  types.Amount `json:"price"`
}

func (d DepositDetails) Encode() string {
	bytes, err := json.Marshal(&d)
	if err != nil {
		panic(errors.Wrap(err, "Failed to encode DepositDetails"))
	}
	return string(bytes)
}
