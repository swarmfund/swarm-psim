package resources

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type WithdrawDetails struct {
	Source string `json:"source"`
	Hash   string `json:"hash"`
}

func (d *WithdrawDetails) Encode() string {
	bytes, err := json.Marshal(&d)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal details"))
	}
	return string(bytes)
}
