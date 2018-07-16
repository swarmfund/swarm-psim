package dashwithdraw

import (
	"testing"

	"gitlab.com/swarmfund/psim/psim/bitcoin"
)

func comparisonHelper(x, y []bitcoin.Out) bool {
	xMap := make(map[bitcoin.Out]int)
	yMap := make(map[bitcoin.Out]int)

	for _, xElem := range x {
		xMap[xElem]++
	}
	for _, yElem := range y {
		yMap[yElem]++
	}

	for xMapKey, xMapVal := range xMap {
		if yMap[xMapKey] != xMapVal {
			return false
		}
	}
	return true
}

func TestGreedyCoinSelector_Fund(t *testing.T) {
	tests := map[string]struct {
		UTXOs          []UTXO
		DustThreshold  int64
		Expected       []bitcoin.Out
		ExpectedChange int64
		Amount         int64
		Error          error
	}{
		"single utxo": { //Amount:5, Change:0, Option:5, Exp:5, Dust: 0
			Amount:         5,
			ExpectedChange: 0,
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

			Expected: []bitcoin.Out{
				{
					Vout:   0,
					TXHash: "hash0",
				},
			},
		},
		"double utxo": { //Amount:2, Change:0, Option:1, 1; Exp:1, 1; Dust: 0
			Amount:         2,
			ExpectedChange: 0,
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
		},
		"big_and_small": { //Amount:3, Change:0, Option:1,2,3; Exp:3; Dust: 0
			Amount:         3,
			ExpectedChange: 0,
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
			Expected: []bitcoin.Out{
				{
					Vout:   2,
					TXHash: "hash2",
				},
			},
		},
		"multiple": { //Amount:15, Change:5, Option:10, 20; Exp:20; Dust: 0
			Amount:         15,
			ExpectedChange: 5,
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
			Expected: []bitcoin.Out{
				{
					Vout:   1,
					TXHash: "hash1",
				},
			},
		},
		"insufficient funds": { //Amount:4, Change:0, Option:1,2; Exp: ; Dust: 0, Err: InsufficientFunds
			Amount:         4,
			ExpectedChange: 0,
			Error:          ErrInsufficientFunds,
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
			Expected: nil,
		},
		"same hash": { //Amount:2, Change:0, Option:1, 1; Exp: 1,1; Dust: 0
			Amount:         2,
			ExpectedChange: 0,
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
		},
		"inactive utxo": { //Amount:3, Change:0, Option: 1,2,3(inactive); Exp: 1,2; Dust: 0
			Amount:         3,
			ExpectedChange: 0,
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
		"dust": { //Amount:20, Change:0, Option: 16,5,2,2,1 Exp:0,1; Dust: 1
			Amount:         20,
			DustThreshold:  1,
			ExpectedChange: 1,
			UTXOs: []UTXO{
				{
					IsInactive: false,
					Value:      16,
					Out: bitcoin.Out{
						Vout:   0,
						TXHash: "hash0",
					},
				},
				{
					IsInactive: false,
					Value:      5,
					Out: bitcoin.Out{
						Vout:   1,
						TXHash: "hash1",
					},
				},
				{
					IsInactive: false,
					Value:      1,
					Out: bitcoin.Out{
						Vout:   2,
						TXHash: "hash2",
					},
				},
				{
					IsInactive: false,
					Value:      2,
					Out: bitcoin.Out{
						Vout:   3,
						TXHash: "hash3",
					},
				},
				{
					IsInactive: false,
					Value:      2,
					Out: bitcoin.Out{
						Vout:   4,
						TXHash: "hash4",
					},
				},
			},
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

		s := NewGreedyCoinSelector(test.DustThreshold)
		for _, u := range test.UTXOs {
			s.AddUTXO(u)
		}

		utxos, change, err := s.Fund(test.Amount)
		if !comparisonHelper(utxos, test.Expected) {
			t.Errorf("%s: expected: %+v\nGot: %+v", k, test.Expected, utxos)
		}
		if change != test.ExpectedChange {
			t.Errorf("%s: expected: %d\nGot: %d", k, test.ExpectedChange, change)
		}
		if err != test.Error {
			t.Errorf("%s: expected: %s\nGot: %s", k, test.Error, err)
		}
	}
}
