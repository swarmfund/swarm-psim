package addrstate

import (
	"sync"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Price struct {
	UpdatedAt time.Time
	Value     int64
}

type externalState struct {
	State     ExternalAccountBindingState
	UpdatedAt time.Time
	Address   string
}

type State struct {
	*sync.Mutex
	//prices   []Price
	//addrs	 sync.Map

	// address -> asset -> balance
	balances map[string]map[string]string
	// external type -> external data -> []events
	external map[int32]map[string][]externalState
	// helper variable for reverse find on external
	internalExternal map[int32]map[string]string
}

func newState() *State {
	return &State{
		Mutex:            &sync.Mutex{},
		external:         map[int32]map[string][]externalState{},
		internalExternal: map[int32]map[string]string{},
		balances:         map[string]map[string]string{},
	}
}

func (s *State) Mutate(ts time.Time, update StateUpdate) {
	s.Lock()
	defer s.Unlock()

	if update.ExternalAccount != nil {
		data := update.ExternalAccount
		switch update.ExternalAccount.State {
		case ExternalAccountBindingStateCreated:
			externalType := data.ExternalType
			if _, ok := s.external[externalType]; !ok {
				s.external[externalType] = map[string][]externalState{}
			}
			s.external[externalType][data.Data] = append(s.external[externalType][data.Data], externalState{
				State:   data.State,
				Address: data.Address,
			})
			if _, ok := s.internalExternal[externalType]; !ok {
				s.internalExternal[externalType] = map[string]string{}
			}
			s.internalExternal[externalType][data.Address] = data.Data
		case ExternalAccountBindingStateDeleted:
			externalType := update.ExternalAccount.ExternalType
			invalidStateErr := errors.From(errors.New("invalid state"), logan.F{
				"reason":        "binding expected to exist",
				"address":       data.Address,
				"external_type": externalType,
			})

			if _, ok := s.internalExternal[externalType]; !ok {
				panic(invalidStateErr)
			}

			external, ok := s.internalExternal[externalType][data.Address]
			if !ok {
				panic(invalidStateErr)
			}

			s.external[externalType][external] = append(s.external[externalType][external], externalState{
				State:   data.State,
				Address: data.Address,
			})
			delete(s.internalExternal[externalType], data.Address)
		default:
			panic(errors.From(errors.New("unknown external update state"), logan.F{
				"external_state": data.State,
			}))
		}
	}

	//if update.AssetPrice != nil {
	//	s.prices = append([]Price{{
	//		UpdatedAt: ts,
	//		Value:     *update.AssetPrice,
	//	}}, s.prices...)
	//}
	//if update.Address != nil {
	//	s.addrs.Store(update.Address.Offchain, update.Address.Tokend)
	//}

	if update.Balance != nil {
		addr := update.Balance.Address
		if s.balances[addr] == nil {
			s.balances[addr] = map[string]string{}
		}
		s.balances[addr][update.Balance.Asset] = update.Balance.Balance
	}
}
