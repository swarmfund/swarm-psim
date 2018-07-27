package dashwithdraw

import (
	"sync"

	"sort"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient funds.")
)

func getKeys(utxos map[bitcoin.Out]UTXO) []bitcoin.Out {
	keys := make([]bitcoin.Out, 0, len(utxos))
	for k := range utxos {
		keys = append(keys, k)
	}

	return keys
}

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

	activeUTXOS := s.getActiveUTXOs()
	if len(activeUTXOS) == 0 {
		// Just in case
		return nil, 0, ErrInsufficientFunds
	}
	su := NewSortedUTXOs(activeUTXOS)

	result, totalFunded := s.chooseUTXOs(su, amount)
	return result, totalFunded - amount, nil
}

func (s GreedyCoinSelector) chooseUTXOs(su *SortedUTXOs, amountToFill int64) ([]bitcoin.Out, int64) {
	for k, v := range su.m {
		if v.Value >= amountToFill && v.Value <= (amountToFill+s.dustThreshold) {
			// Ideal UTXO found
			result := make([]bitcoin.Out, 0)
			result = append(result, k)
			return result, v.Value
		}
	}

	//Ideal UTXO not found, choosing from 2 groups
	biggerUTXOS, smallerUTXOS := s.splitAndChoose(su, amountToFill)

	//If small UTXOs can fulfil amount
	if smallerUTXOS != nil {
		return s.greedyKnapsack(smallerUTXOS, amountToFill)
	}

	//choosing one minimal UTXO from the biggerUTXOS
	chosen, utxo, err := biggerUTXOS.PopSmallest()
	if err != nil {
		return nil, 0
	}

	result := make([]bitcoin.Out, 0)
	result = append(result, chosen)
	return result, utxo.Value
}

func (s GreedyCoinSelector) greedyKnapsack(su *SortedUTXOs, amountToFill int64) (chosen []bitcoin.Out, filled int64) {
	if su.Len() == 0 {
		return nil, 0
	}

	var sum int64
	chosen = make([]bitcoin.Out, 0)
	for sum < amountToFill {
		k, v, err := su.PopBiggest()
		if err != nil {
			return nil, 0
		}
		sum += v.Value
		chosen = append(chosen, k)
	}
	return chosen, sum
}

//splitAndChoose splits UTXOs into 2 groups -
//the ones that bigger than amountToFill,
//and the ones that smaller than amountToFill
//Returns only one map - smallerUTXOS if there are sufficient sum
//or biggerUTXOS if smallerUTXOS are not enough to fulfill amountToFill
func (s GreedyCoinSelector) splitAndChoose(su *SortedUTXOs, amountToFill int64) (biggerUTXOS *SortedUTXOs, smallerUTXOS *SortedUTXOs) {
	var totalSmalls int64
	smaller := make(map[bitcoin.Out]UTXO)
	bigger := make(map[bitcoin.Out]UTXO)
	for k, v := range su.m {
		if v.Value < amountToFill {
			smaller[k] = v
			totalSmalls += v.Value
			continue
		}
		bigger[k] = v
	}

	if totalSmalls < amountToFill {
		biggerUTXOS = NewSortedUTXOs(bigger)
		sort.Sort(biggerUTXOS)
		return biggerUTXOS, nil
	}

	smallerUTXOS = NewSortedUTXOs(smaller)
	sort.Sort(smallerUTXOS)
	return nil, smallerUTXOS
}

func (s GreedyCoinSelector) getActiveUTXOs() map[bitcoin.Out]UTXO {
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
