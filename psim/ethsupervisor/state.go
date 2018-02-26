package ethsupervisor

import (
	"time"
	"context"
)

type State interface {
	AddressAt(context.Context, time.Time, string) *string
	Balance(context.Context, string) *string
}
