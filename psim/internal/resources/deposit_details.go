package resources

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/horizon-connector/v2/types"
)

type DepositDetails struct {
	TXHash string       `json:"tx_hash"`
	Price  types.Amount `json:"price"`
}

func (d DepositDetails) Encode() string {
	bytes, err := json.Marshal(&d)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal details"))
	}
	return string(bytes)
}
