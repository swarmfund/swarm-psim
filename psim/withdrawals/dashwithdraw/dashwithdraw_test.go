package dashwithdraw

import (
	"testing"

	"gitlab.com/swarmfund/psim/psim/bitcoin"

	//"github.com/magiconair/properties/assert"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/assert"
)

func TestRandomCoinSelector_Fund(t *testing.T) {
	tests := map[string]struct {
		UTXOs    []UTXO
		Expected []bitcoin.Out
		Change   int64
		Amount   int64
	}{
		"single utxo": {
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      500000000,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
			},
			Change: 0,
			Expected: []bitcoin.Out{
				{
					Vout:   0,
					TXHash: "hash0",
				},
			},
			Amount: 500000000,
		},
		"double utxo": {
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      100000000,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      100000000,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
			},
			Change: 0,
			Expected: []bitcoin.Out{
				{
					Vout:   1,
					TXHash: "hash1",
				},
				{
					Vout:   0,
					TXHash: "hash0",
				},
			},
			Amount: 200000000,
		},
		"big_and_small": {
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      100000000,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      200000000,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
				{
					IsInactive: false,
					Value:      300000000,
					Out: bitcoin.Out{
						Vout:   2,
						TXHash: "hash2",
					},
				},
			},
			Change: 0,
			Expected: []bitcoin.Out{
				{
					Vout:   2,
					TXHash: "hash2",
				},
			},
			Amount: 300000000,
		},
		"multiple": {
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      100000000,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      200000000,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
			},
			Change: 50000000,
			Expected: []bitcoin.Out{
				{
					Vout:   1,
					TXHash: "hash1",
				},
			},
			Amount: 150000000,
		},
	}

	for k, test := range tests {
		println(k)

		s := NewRandomCoinSelector(0)
		for _, u := range test.UTXOs {
			s.AddUTXO(u)
		}

		utxos, change, _ := s.Fund(test.Amount)
		assert.ElementsMatch(t, utxos, test.Expected)
		assert.EqualValues(t, test.Change, change)
	}
}
