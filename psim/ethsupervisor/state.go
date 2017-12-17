package ethsupervisor

import (
	"time"
)

type State interface {
	AddressAt(time.Time, string) *string
	PriceAt(time.Time) *int64
	Balance(string) *string
}
