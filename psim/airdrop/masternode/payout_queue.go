package masternode

import "time"

type PayoutQueue struct {
	FirstPayout   time.Time
	BlockDuration time.Duration
}
