package dashwithdraw

import "fmt"

type ErrTooBigFeePerKB struct {
	suggestedFee  float64
	maxAllowedFee float64
}

func NewErrTooBigFeePerKB(suggestedFee, maxAllowedFee float64) ErrTooBigFeePerKB {
	return ErrTooBigFeePerKB{
		suggestedFee:  suggestedFee,
		maxAllowedFee: maxAllowedFee,
	}
}

func (e ErrTooBigFeePerKB) Error() string {
	return fmt.Sprintf("Suggested fee per KB (%f) is too big, max allowed is (%f).", e.suggestedFee, e.maxAllowedFee)
}
