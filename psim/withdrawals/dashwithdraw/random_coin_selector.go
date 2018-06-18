package dashwithdraw

import (
	"sync"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

type RandomCoinSelector struct {
	mu sync.Mutex

	utxos map[bitcoin.Out]UTXO
}

func NewRandomCoinSelector() *RandomCoinSelector {
	return &RandomCoinSelector{
		mu: sync.Mutex{},

		utxos: make(map[bitcoin.Out]UTXO),
	}
}

func (s RandomCoinSelector) AddUTXO(utxo UTXO) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.utxos[utxo.Out] = utxo
}

// TODO Use mutex
// TODO
func (s RandomCoinSelector) Fund(amount int64) (utxos []bitcoin.Out, change int64, err error) {
	// TODO
	return nil, 0, errors.New("Not implemented.")
}

func (s RandomCoinSelector) TryRemoveUTXO(out bitcoin.Out) bool {
	utxo, ok := s.utxos[out]
	if !ok {
		// Does not exist
		return false
	}

	// Exists
	utxo.IsInactive = true

	s.mu.Lock()
	defer s.mu.Unlock()
	s.utxos[out] = utxo

	return false
}

func (s RandomCoinSelector) getActiveBalance() int64 {
	var activeBalance int64

	for _, utxo := range s.utxos {
		if !utxo.IsInactive {
			// UTXO is active
			activeBalance += utxo.Value
		}
	}

	return activeBalance
}
