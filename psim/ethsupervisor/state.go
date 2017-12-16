package ethsupervisor

import (
	"context"
	"time"
)

type State interface {
	AddressAt(context.Context, time.Time, string) *string
	PriceAt(context.Context, time.Time) *int64
}
