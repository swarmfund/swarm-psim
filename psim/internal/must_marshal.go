package internal

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MustMarshal(a interface{}) string {
	bb, err := json.Marshal(a)
	if err != nil {
		panic(errors.Wrap(err, "failed to marshal", logan.F{
			"the_value": a,
		}))
	}

	return string(bb)
}
