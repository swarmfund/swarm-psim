package dashwithdraw

import (
	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

type SortableUTXOs struct {
	m    map[bitcoin.Out]UTXO
	keys []bitcoin.Out
}

func NewSortableUTXOs(utxos map[bitcoin.Out]UTXO) *SortableUTXOs {
	return &SortableUTXOs{m: utxos,
		keys: getKeys(utxos),
	}
}

func (su *SortableUTXOs) Len() int {
	return len(su.m)
}

//Less doing exact opposite to sort in descending order
func (su *SortableUTXOs) Less(i, j int) bool {
	return su.m[su.keys[i]].Value > su.m[su.keys[j]].Value
}

func (su *SortableUTXOs) Swap(i, j int) {
	su.keys[i], su.keys[j] = su.keys[j], su.keys[i]
}

func (su *SortableUTXOs) PopBiggest() (key bitcoin.Out, value UTXO) {
	if len(su.keys) == 0 {
		panic("Container is empty")
	}
	key = su.keys[0]
	value = su.m[key]
	su.keys = append(su.keys[:0], su.keys[1:]...)
	delete(su.m, key)
	return key, value
}
