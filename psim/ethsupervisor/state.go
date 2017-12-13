package ethsupervisor

import "time"

type State interface {
	AddressAt(time.Time, string) *string
}
