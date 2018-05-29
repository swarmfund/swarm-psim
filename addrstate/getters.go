package addrstate

import (
	"context"
	"time"
)

func (w *Watcher) ExternalAccountAt(ctx context.Context, ts time.Time, systemType int32, data string) *string {
	w.ensureReached(ctx, ts)

	w.state.Lock()
	defer w.state.Unlock()

	if _, ok := w.state.external[systemType]; !ok {
		// external system states doesn't exist (yet)
		return nil
	}

	states := w.state.external[systemType][data]
	if len(states) == 0 {
		return nil
	}
	// iterating through the closed periods
	for i := 0; i < len(states)-1; i += 1 {
		a := states[i]
		b := states[i+1]
		if ts.After(a.UpdatedAt) && ts.Before(b.UpdatedAt) {
			// we found time interval that includes our ts,
			// first states is current one
			addr := a.Address
			return &addr
		}
	}
	// checking last known state
	lastState := states[len(states)-1]
	if ts.After(lastState.UpdatedAt) && lastState.State == ExternalAccountBindingStateCreated {
		addr := lastState.Address
		return &addr
	}
	// seems like rogue deposit, but who cares
	return nil
}

// BindExternalSystemEntities returns all known external data for systemType
func (w *Watcher) BindedExternalSystemEntities(ctx context.Context, systemType int32) (result []string) {
	w.ensureReached(ctx, time.Now())
	w.state.Lock()
	defer w.state.Unlock()

	if _, ok := w.state.external[systemType]; !ok {
		return result
	}

	entities := w.state.external[systemType]
	for entity, _ := range entities {
		result = append(result, entity)
	}
	return result
}

func (w *Watcher) Balance(ctx context.Context, address string, asset string) *string {
	w.state.Lock()
	defer w.state.Unlock()

	// let's hope for the best and try to get balance before reaching head
	if w.state.balances[address] != nil {
		if balance, ok := w.state.balances[address][asset]; ok {
			return &balance
		}
	}

	// if we don't have balance already, let's wait for latest ledger
	w.ensureReached(ctx, time.Now())

	// now check again
	if w.state.balances[address] != nil {
		if balance, ok := w.state.balances[address][asset]; ok {
			return &balance
		}
	}

	// seems like user doesn't have balance in provided asset atm
	return nil
}
