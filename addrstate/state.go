package addrstate

import (
	"time"
	"sync"
)

type Price struct {
	UpdatedAt time.Time
	Value     int64
}

type State struct {
	prices   []Price
	addrs	 sync.Map
	balances sync.Map
}

func newState() *State {
	return &State{
		addrs:    sync.Map{},
		balances: sync.Map{},
	}
}

func (s *State) Mutate(ts time.Time, update StateUpdate) {
	if update.AssetPrice != nil {
		s.prices = append([]Price{{
			UpdatedAt: ts,
			Value:     *update.AssetPrice,
		}}, s.prices...)
	}
	if update.Address != nil {
		s.addrs.Store(update.Address.Offchain, update.Address.Tokend)
	}
	if update.Balance != nil {
		s.balances.Store(update.Balance.Address, update.Balance.Balance)
	}
}
