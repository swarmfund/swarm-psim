package mrefairdrop

import (
	"testing"

	"gitlab.com/tokend/go/amount"
)

func Test_countIssuanceAmount(t *testing.T) {
	table := []struct {
		balance uint64
		result  uint64
	}{
		{amount.One * 100, amount.One * 20},
		{amount.One * 5, amount.One * 1},
		{amount.One * 500, amount.One * 100},
		{amount.One * 100000, amount.One * 20000},
		{amount.One * 1000000, amount.One * 20000},
	}

	for _, tc := range table {
		got := countIssuanceAmount(tc.balance)
		if got != tc.result {
			t.Errorf("got %d, want %d", got, tc.result)
		}
	}
}
