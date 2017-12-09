package resource

import (
	"encoding/json"

	. "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/swarmfund/go/amount"
)

const minStripeAmountValue = 5000 // 50 cents

// StripeChargeRequest is the struct to unmarshal the request of Charge creation into.
type StripeChargeRequest struct {
	Token     string `json:"token"`
	Amount    int64  `json:"amount"`
	Reference string `json:"reference"`
	Receiver  string `json:"receiver"`
	Asset     string `json:"asset"`
}

// Validate implements interface ape.Validator.
func (r *StripeChargeRequest) Validate() error {
	return ValidateStruct(r,
		Field(&r.Token, Required),
		Field(&r.Amount, Required, Min(minStripeAmountValue)),
		Field(&r.Reference, Required),
	)
}

func (r *StripeChargeRequest) UnmarshalJSON(data []byte) error {
	type request StripeChargeRequest
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
		r.Amount, err = amount.Parse(rr.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}
