package dashwithdraw

import (
	"testing"

	"gitlab.com/swarmfund/psim/psim/bitcoin"
	"reflect"
)

func TestRandomCoinSelector_Fund(t *testing.T) {
	tests := map[string]struct {
		UTXOs          []UTXO
		Expected       []bitcoin.Out
		ExpectedChange int64
		Amount         int64
		Error          error
	}{
		"single utxo": {
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      5,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
			},
			ExpectedChange: 0,
			Expected: []bitcoin.Out{
				{
					Vout:   0,
					TXHash: "hash0",
				},
			},
			Amount: 5,
		},
		"double utxo": {
			Amount: 2,
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
			},
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
			ExpectedChange: 0,

		},
		"big_and_small": {
			Amount: 3,
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      2,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
				{
					IsInactive: false,
					Value:      3,
					Out: bitcoin.Out{
						Vout:   2,
						TXHash: "hash2",
					},
				},
			},
			ExpectedChange: 0,
			Expected: []bitcoin.Out{
				{
					Vout:   2,
					TXHash: "hash2",
				},
			},
		},
		"multiple": {
			Amount: 15,
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      10,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      20,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
			},
			ExpectedChange: 5,
			Expected: []bitcoin.Out{
				{
					Vout:   1,
					TXHash: "hash1",
				},
			},
		},
		"insufficient funds":{
			Amount: 4,
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      2,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
			},
			ExpectedChange: 0,
			Expected: nil,
			Error: ErrInsufficientFunds,
		},
		"same hash": {
			Amount: 2,
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash0",
					},
				},
			},
			Expected: []bitcoin.Out{
				{
					Vout:   1,
					TXHash: "hash0",
				},
				{
					Vout:   0,
					TXHash: "hash0",
				},
			},
			ExpectedChange: 0,

		},
		"inactive utxo": {
			Amount: 3,
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      2,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
				{
					IsInactive: true,
					Value:      3,
					Out: bitcoin.Out{
						Vout:   2,
						TXHash: "hash2",
					},
				},
			},
			ExpectedChange: 0,
			Expected: []bitcoin.Out{
				{
					Vout:   0,
					TXHash: "hash0",
				},
				{
					Vout:   1,
					TXHash: "hash1",
				},
			},
		},
	}

	for k, test := range tests {

		s := NewRandomCoinSelector(0)
		for _, u := range test.UTXOs {
			s.AddUTXO(u)
		}

		utxos, change, _ := s.Fund(test.Amount)
		if !reflect.DeepEqual(utxos, test.Expected){
			t.Errorf("%s: expected: ")
		}
	}
}
