package dashwithdraw

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

type SortedUTXOs struct {
	m    map[bitcoin.Out]UTXO
	keys []bitcoin.Out
}

func NewSortedUTXOs(utxos map[bitcoin.Out]UTXO) *SortedUTXOs {
	return &SortedUTXOs{m: utxos,
		keys: getKeys(utxos),
	}
}

func (su *SortedUTXOs) Len() int {
	return len(su.m)
}

//Less doing exact opposite to sort in descending order
func (su *SortedUTXOs) Less(i, j int) bool {
	return su.m[su.keys[i]].Value > su.m[su.keys[j]].Value
}

func (su *SortedUTXOs) Swap(i, j int) {
	su.keys[i], su.keys[j] = su.keys[j], su.keys[i]
}

func (su *SortedUTXOs) PopBiggest() (key bitcoin.Out, value UTXO, err error) {
	if len(su.keys) == 0 {
		return bitcoin.Out{}, UTXO{}, errors.New("Container is empty")
	}
	key = su.keys[0]
	value = su.m[key]
	su.keys = append(su.keys[:0], su.keys[1:]...)
	delete(su.m, key)
	return key, value, nil
}

func (su *SortedUTXOs) PopSmallest() (key bitcoin.Out, value UTXO, err error) {
	if len(su.keys) == 0 {
		return bitcoin.Out{}, UTXO{}, errors.New("Container is empty")
	}
	key = su.keys[len(su.keys)-1]
	value = su.m[key]
	su.keys = append(su.keys[:0], su.keys[:len(su.keys)-1]...)
	delete(su.m, key)
	return key, value, nil
}
