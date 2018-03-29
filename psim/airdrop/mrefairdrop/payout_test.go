package mrefairdrop

import (
	"testing"

	"gitlab.com/swarmfund/go/amount"
)

func Test_countIssuanceAmount(t *testing.T) {
	table := []struct {
		referrals int
		balance   uint64
		result    uint64
	}{
		{15, amount.One * 700, amount.One * 180},
		{11, amount.One * 400000, amount.One * 4055},
		{20, amount.One * 400000, amount.One * 4100},
		{1000, amount.One * 10000, amount.One * 9000},
		{1000, amount.One * 100, amount.One * 6000},
		{1000, amount.One * 0, amount.One * 5000},
		{2000, amount.One * 10000, amount.One * 14000},
		{3000, amount.One * 10000, amount.One * 19000},
		{8000, amount.One * 10000, amount.One * 20000},
		{8000, amount.One * 0, amount.One * 20000},
	}

	for _, tc := range table {
		got := countIssuanceAmount(tc.referrals, tc.balance)
		if got != tc.result {
			t.Errorf("got %d, want %d", got, tc.result)
		}
	}
}
