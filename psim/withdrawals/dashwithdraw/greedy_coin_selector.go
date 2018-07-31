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
	result, totalFunded := s.chooseUTXOs(activeUTXOS, amount)
	return result, totalFunded - amount, nil
}

func (s GreedyCoinSelector) chooseUTXOs(utxos map[bitcoin.Out]UTXO, amountToFill int64) ([]bitcoin.Out, int64) {
	for k, v := range utxos {
		if v.Value >= amountToFill && v.Value <= (amountToFill+s.dustThreshold) {
			// Ideal UTXO found
			return []bitcoin.Out{k}, v.Value
		}
	}

	//Ideal UTXO not found, choosing from 2 groups
	biggerUTXOS, smallerUTXOS := s.splitAndChoose(utxos, amountToFill)

	//If small UTXOs can fulfil amount
	if smallerUTXOS != nil {
		return s.greedyKnapsack(smallerUTXOS, amountToFill)
	}

	//choosing one minimal UTXO from the biggerUTXOS
	chosen, utxo := smallestInMap(biggerUTXOS)
	return []bitcoin.Out{chosen}, utxo.Value
}

func (s GreedyCoinSelector) greedyKnapsack(utxos map[bitcoin.Out]UTXO, amountToFill int64) (chosen []bitcoin.Out, filled int64) {
	if len(utxos) == 0 {
		return nil, 0
	}
	su := NewSortableUTXOs(utxos)
	sort.Sort(su)

	var sum int64
	chosen = make([]bitcoin.Out, 0)
	for sum < amountToFill {
		k, v := su.PopBiggest()
		sum += v.Value
		chosen = append(chosen, k)
	}
	return chosen, sum
}

//SplitAndChoose splits UTXOs into 2 groups -
//the ones that bigger than amountToFill,
//and the ones that smaller than amountToFill
//Returns only one map - smallerUTXOS if there are sufficient sum
//or biggerUTXOS if smallerUTXOS are not enough to fulfill amountToFill
func (s GreedyCoinSelector) splitAndChoose(utxos map[bitcoin.Out]UTXO, amountToFill int64) (bigger map[bitcoin.Out]UTXO, smaller map[bitcoin.Out]UTXO) {
	var totalSmalls int64
	smaller = make(map[bitcoin.Out]UTXO)
	bigger = make(map[bitcoin.Out]UTXO)
	for k, v := range utxos {
		if v.Value < amountToFill {
			smaller[k] = v
			totalSmalls += v.Value
			continue
		}
		bigger[k] = v
	}

	if totalSmalls < amountToFill {
		return bigger, nil
	}

	return nil, smaller
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

func smallestInMap(utxos map[bitcoin.Out]UTXO) (bitcoin.Out, UTXO) {
	keys := getKeys(utxos)
	min := keys[0]
	for _, k := range keys {
		if utxos[k].Value < utxos[min].Value {
			min = k
		}
	}
	return min, utxos[min]
}
