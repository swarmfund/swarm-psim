package dashwithdraw

import (
	"sync"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient funds.")
)

type RandomCoinSelector struct {
	dustThreshold int64
	mu            sync.Mutex
	utxos         map[bitcoin.Out]UTXO
}

func NewRandomCoinSelector(dustThreshold int64) *RandomCoinSelector {
	return &RandomCoinSelector{
		dustThreshold: dustThreshold,

		mu:    sync.Mutex{},
		utxos: make(map[bitcoin.Out]UTXO),
	}
}

func (s RandomCoinSelector) AddUTXO(utxo UTXO) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.utxos[utxo.Out] = utxo
}

func (s RandomCoinSelector) Fund(amount int64) (utxos []bitcoin.Out, change int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	activeBalance := s.getActiveBalance()
	if activeBalance < amount {
		return nil, 0, ErrInsufficientFunds
	}

	activeUTXOSCopy := make(map[bitcoin.Out]UTXO)
	for k, v := range s.utxos {
		if !v.IsInactive {
			activeUTXOSCopy[k] = v
		}
	}

	var totalFunded int64
	var result []bitcoin.Out

	for totalFunded < amount {
		if len(activeUTXOSCopy) == 0 {
			// Just in case
			return nil, 0, ErrInsufficientFunds
		}

		chosenOut, chosenUTXO := s.chooseUTXO(activeUTXOSCopy, amount-totalFunded)
		result = append(result, chosenOut)
		totalFunded += chosenUTXO.Value
		delete(activeUTXOSCopy, chosenOut)
	}

	return result, totalFunded - amount, nil
}

func (s RandomCoinSelector) chooseUTXO(utxos map[bitcoin.Out]UTXO, amountToFill int64) (bitcoin.Out, UTXO) {
	for k, v := range utxos {
		if v.Value >= amountToFill && v.Value <= (amountToFill+s.dustThreshold) {
			// Ideal UTXO found
			return k, v
		}
	}

	var desired bitcoin.Out
	for k, v := range utxos {
		if v.Value > utxos[desired].Value {
			desired = k
		}
	}

	return desired, utxos[desired]
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
