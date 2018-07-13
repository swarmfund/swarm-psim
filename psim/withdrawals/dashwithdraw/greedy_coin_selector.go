package dashwithdraw

import (
	"sync"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient funds.")
)

const (
	greedy = iota
	biggest = iota
)

type GreedyCoinSelector struct {
	dustThreshold int64
	mu            sync.Mutex
	utxos         map[bitcoin.Out]UTXO
}

func NewGreedyCoinSelector(dustThreshold int64) *GreedyCoinSelector {
	return &GreedyCoinSelector{
		dustThreshold: dustThreshold,

		mu:    sync.Mutex{},
		utxos: make(map[bitcoin.Out]UTXO),
	}
}

func (s GreedyCoinSelector) AddUTXO(utxo UTXO) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.utxos[utxo.Out] = utxo
}

func (s GreedyCoinSelector) Fund(amount int64) (utxos []bitcoin.Out, change int64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	activeBalance := s.getActiveBalance()
	if activeBalance < amount {
		return nil, 0, ErrInsufficientFunds
	}

	var totalFunded int64
	var result []bitcoin.Out
	mode := greedy
	activeUTXOS := s.getActiveUTXOs()
	for totalFunded < amount {
		if len(activeUTXOS) == 0 {
			// Just in case
			return nil, 0, ErrInsufficientFunds
		}

		chosenOut, chosenUTXO := s.chooseUTXO(activeUTXOS, amount-totalFunded, mode)
		//if went wrong way
		if _, ok := activeUTXOS[chosenOut]; !ok{

			mode = biggest
			activeUTXOS = s.getActiveUTXOs()
			result = make([]bitcoin.Out, 0)
			totalFunded = 0
			continue
		}
		result = append(result, chosenOut)
		totalFunded += chosenUTXO.Value
		delete(activeUTXOS, chosenOut)
	}

	return result, totalFunded - amount, nil
}

func (s GreedyCoinSelector) chooseUTXO(utxos map[bitcoin.Out]UTXO, amountToFill int64, mode int) (bitcoin.Out, UTXO) {
	var desired bitcoin.Out
	switch  mode {
	case greedy:
		for k, v := range utxos {
			if v.Value > utxos[desired].Value && v.Value < amountToFill {
				desired = k
			}
		}
	case biggest:
		for k, v := range utxos {
			if v.Value > utxos[desired].Value{
				desired = k
			}
		}
	}
	return desired, utxos[desired]
}

func (s GreedyCoinSelector) getActiveUTXOs() map[bitcoin.Out]UTXO{
	activeUTXOSCopy := make(map[bitcoin.Out]UTXO)
	for k, v := range s.utxos {
		if !v.IsInactive {
			activeUTXOSCopy[k] = v
		}
	}
	return activeUTXOSCopy
}


func (s GreedyCoinSelector) TryRemoveUTXO(out bitcoin.Out) bool {
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

func (s GreedyCoinSelector) getActiveBalance() int64 {
	var activeBalance int64

	for _, utxo := range s.utxos {
		if !utxo.IsInactive {
			// UTXO is active
			activeBalance += utxo.Value
		}
	}

	return activeBalance
}
