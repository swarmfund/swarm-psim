package ethsupervisor

import (
	"time"
	"context"
)

type State interface {
	AddressAt(context.Context, time.Time, string) *string
	PriceAt(context.Context, time.Time) *int64
	Balance(context.Context, string) *string
}
