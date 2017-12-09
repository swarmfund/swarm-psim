package resource

import (
	"encoding/json"
	"strconv"

	. "github.com/go-ozzo/ozzo-validation"
)

// TODO Move this to some config
const minSatoshiAmountValue uint64 = 1000000 // 0.01 BTC

// BitcoinChargeRequest is the struct to unmarshal the request of Charge creation into.
type BitcoinChargeRequest struct {
	Amount    uint64 `json:"amount"` // In satoshi
	Receiver  string `json:"receiver"` // Should be Stellar AccountID
}

// Validate implements interface ape.Validator.
func (r *BitcoinChargeRequest) Validate() error {
	return ValidateStruct(r,
		Field(&r.Amount, Required, Min(minSatoshiAmountValue)),
		Field(&r.Receiver, Required),
	)
}

func (r *BitcoinChargeRequest) UnmarshalJSON(data []byte) error {
	type request BitcoinChargeRequest
	rr := &struct {
		Amount string `json:"amount"`
		*request
	}{
		request: (*request)(r),
	}

	err := json.Unmarshal(data, &rr)
	if err != nil {
		return err
	}
	if rr.Amount != "" {
		r.Amount, err = strconv.ParseUint(rr.Amount, 10, 64)
		if err != nil {
			return err
		}
	}
	return nil
}
