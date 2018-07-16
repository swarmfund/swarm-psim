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

	//Ideal UTXO not found, splitting into 2 groups
	smallerUTXOS := make(map[bitcoin.Out]UTXO)
	biggerUTXOS := make(map[bitcoin.Out]UTXO)
	var totalSmalls int64

	//Split UTXOs into 2 groups - the ones that bigger than amountToFill,
	//and the ones that smaller than amountToFill
	//also we check whether sum of all small utxos is greater or equal to amountToFill
	for k, v := range utxos {
		if v.Value < amountToFill {
			smallerUTXOS[k] = v
			totalSmalls += v.Value
			continue
		}
		biggerUTXOS[k] = v
	}

	//If small UTXOs can fulfil amount
	if totalSmalls >= amountToFill {
		keys := getKeys(smallerUTXOS)
		permutations, isEmpty := s.getPermutations(smallerUTXOS, keys, amountToFill)

		if !isEmpty {
			indexToChoose := rand.Intn(len(permutations))
			sum := sumOfUTXOs(permutations[indexToChoose])
			for i, p := range permutations {
				newSum := sumOfUTXOs(p)
				if newSum < sum {
					sum = newSum
					indexToChoose = i
				}
			}

			return permutations[indexToChoose], sum
		}
	}

	//choosing one minimal UTXO from the ones that biggerUTXOS than amount
	minOut := chooseRandomInMap(biggerUTXOS)
	for i, utxo := range biggerUTXOS {
		if utxo.Value < biggerUTXOS[minOut].Value {
			minOut = i
		}
	}
	result := make(map[bitcoin.Out]UTXO)
	result[minOut] = biggerUTXOS[minOut]

	return result, result[minOut].Value
}

func sumOfUTXOs(utxos map[bitcoin.Out]UTXO) int64 {
	var sum int64
	for _, v := range utxos {
		sum += v.Value
	}
	return sum
}

func chooseRandomInMap(utxos map[bitcoin.Out]UTXO) bitcoin.Out {
	indexToChoose := rand.Intn(len(utxos))
	var i int
	for k := range utxos {
		if i == indexToChoose {
			return k
		}
		i++
	}
	return bitcoin.Out{}
}

func (s GreedyCoinSelector) getPermutations(utxos map[bitcoin.Out]UTXO, keys []bitcoin.Out, amountToFill int64) ([]map[bitcoin.Out]UTXO, bool) {
	rounds := min(1<<uint(len(utxos)), 1000)

	options := powerSet(rounds, keys)

	result := make([]map[bitcoin.Out]UTXO, 0, rounds)

	for _, option := range options {

		combination := make(map[bitcoin.Out]UTXO)
		var sum int64
		for _, key := range option {
			combination[key] = utxos[key]
			sum += utxos[key].Value
		}
		if sum < amountToFill {
			continue
		}
		result = append(result, combination)
	}

	var isEmpty bool
	if len(result) == 0 {
		isEmpty = true
	}
	return result, isEmpty
}

func getKeys(utxos map[bitcoin.Out]UTXO) []bitcoin.Out {
	keys := make([]bitcoin.Out, 0, len(utxos))
	for k := range utxos {
		keys = append(keys, k)
	}

	return keys
}
func powerSet(powerSetSize int, original []bitcoin.Out) [][]bitcoin.Out {
	result := make([][]bitcoin.Out, 0, powerSetSize)

	var index int
	for index < powerSetSize {
		var subSet []bitcoin.Out

		for j, elem := range original {
			if index&(1<<uint(j)) > 0 {
				subSet = append(subSet, elem)
			}
		}
		result = append(result, subSet)
		index++
	}
	return result
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
