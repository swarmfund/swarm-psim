package dashwithdraw

import (
	"sync"

	"math/rand"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

var (
	ErrInsufficientFunds = errors.New("Insufficient funds.")
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

	activeUTXOS := s.getActiveUTXOs()

	if len(activeUTXOS) == 0 {
		// Just in case
		return nil, 0, ErrInsufficientFunds
	}

	var totalFunded int64
	var result map[bitcoin.Out]UTXO
	result, totalFunded = s.chooseUTXOs(activeUTXOS, amount)

	return getKeys(result), totalFunded - amount, nil

}

func (s GreedyCoinSelector) chooseUTXOs(utxos map[bitcoin.Out]UTXO, amountToFill int64) (map[bitcoin.Out]UTXO, int64) {
	for k, v := range utxos {
		if v.Value >= amountToFill && v.Value <= (amountToFill+s.dustThreshold) {
			// Ideal UTXO found
			result := make(map[bitcoin.Out]UTXO)
			result[k] = v
			return result, v.Value
		}
	}

	//Ideal UTXO not found, choosing from 2 groups
	biggerUTXOS, smallerUTXOS := s.splitAndChoose(utxos, amountToFill)

	//If small UTXOs can fulfil amount
	if smallerUTXOS != nil {
		var sum int64
		chosen := make(map[bitcoin.Out]UTXO)
		for sum < amountToFill {
			max := s.chooseMaxUTXO(smallerUTXOS)
			sum += smallerUTXOS[max].Value

			chosen[max] = smallerUTXOS[max]
			delete(smallerUTXOS, max)

		}
		return chosen, sum
	}

	//choosing one minimal UTXO from the biggerUTXOS
	chosen := s.chooseMinUTXO(biggerUTXOS)
	result := make(map[bitcoin.Out]UTXO)
	result[chosen] = biggerUTXOS[chosen]

	return result, result[chosen].Value
}

//Split UTXOs into 2 groups - the ones that bigger than amountToFill,
//and the ones that smaller than amountToFill
//Returns only one map - smallerUTXOS if there are sufficient sum
//or biggerUTXOS if smallerUTXOS are not enough to fulfill amountToFill
func (s GreedyCoinSelector) splitAndChoose(utxos map[bitcoin.Out]UTXO, amountToFill int64) (map[bitcoin.Out]UTXO, map[bitcoin.Out]UTXO) {
	var totalSmalls int64
	smallerUTXOS := make(map[bitcoin.Out]UTXO)
	biggerUTXOS := make(map[bitcoin.Out]UTXO)
	for k, v := range utxos {
		if v.Value < amountToFill {
			smallerUTXOS[k] = v
			totalSmalls += v.Value
			continue
		}
		biggerUTXOS[k] = v
	}

	if totalSmalls < amountToFill {
		return biggerUTXOS, nil
	}
	return nil, smallerUTXOS

}

func (s GreedyCoinSelector) chooseMaxUTXO(utxos map[bitcoin.Out]UTXO) bitcoin.Out {
	max := s.chooseRandomInMap(utxos)
	for k, v := range utxos {
		if utxos[max].Value < v.Value {
			max = k
		}
	}

	return max
}

func (s GreedyCoinSelector) chooseMinUTXO(utxos map[bitcoin.Out]UTXO) bitcoin.Out {
	min := s.chooseRandomInMap(utxos)
	for k, v := range utxos {
		if utxos[min].Value > v.Value {
			min = k
		}
	}
	return min
}
func (s GreedyCoinSelector) chooseRandomInMap(utxos map[bitcoin.Out]UTXO) bitcoin.Out {
	if len(utxos) > 0 {
		indexToChoose := rand.Intn(len(utxos))
		var i int
		for k := range utxos {
			if i == indexToChoose {
				return k
			}
			i++
		}
	}
	return bitcoin.Out{}
}

func getKeys(utxos map[bitcoin.Out]UTXO) []bitcoin.Out {
	keys := make([]bitcoin.Out, 0, len(utxos))
	for k := range utxos {
		keys = append(keys, k)
	}

	return keys
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
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
