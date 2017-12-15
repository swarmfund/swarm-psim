package addrstate

import "time"

type Price struct {
	UpdatedAt time.Time
	Value     int64
}

type State struct {
	prices []Price
	addrs  map[string]string
}

func newState() *State {
	return &State{
		addrs: map[string]string{},
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
		s.addrs[update.Address.Offchain] = update.Address.Tokend
	}
}
