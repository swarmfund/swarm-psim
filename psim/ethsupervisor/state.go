package ethsupervisor

import (
	"context"
	"time"
)

type State interface {
	ExternalAccountAt(context.Context, time.Time, int32, string) *string
	Balance(context.Context, string, string) *string
}
