package ethsupervisor

import (
	"time"
	"context"
)

type State interface {
	AddressAt(context.Context, time.Time, string) *string
}
