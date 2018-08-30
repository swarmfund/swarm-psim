package masternode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func pvalequal(t *testing.T, a string, b *string) {
	t.Helper()
	if b == nil {
		t.Fatalf("expected '%s' got nil", a)
	}
	if a != *b {
		t.Fatalf("expected '%s' got '%s'", a, *b)
	}
}

func TestRewardQueue(t *testing.T) {
	t.Run("empty queue should not return next", func(t *testing.T) {
		rq := RewardQueue{}
		assert.Nil(t, rq.Next())
	})

	t.Run("promote after should be respected", func(t *testing.T) {
		rq := RewardQueue{
			promoteAfter: 2,
		}
		rq.Add("addr")
		assert.Nil(t, rq.Next())
		assert.Nil(t, rq.Next())
		pvalequal(t, "addr", rq.Next())
	})

	t.Run("add should be idempotent", func(t *testing.T) {
		first := "first"
		second := "second"
		third := "third"
		rq := RewardQueue{}
		rq.Add(first)
		rq.Add(first)
		rq.Add(second)
		rq.Add(first)
		rq.Add(first)
		rq.Add(third)
		pvalequal(t, first, rq.Next())
		pvalequal(t, second, rq.Next())
		pvalequal(t, third, rq.Next())
	})

	t.Run("single address should always be next", func(t *testing.T) {
		addr := "addr"
		rq := RewardQueue{}
		rq.Add(addr)
		pvalequal(t, addr, rq.Next())
		pvalequal(t, addr, rq.Next())
		pvalequal(t, addr, rq.Next())

	})

	t.Run("multiple addresses should round-robin", func(t *testing.T) {
		first := "first"
		second := "second"
		rq := RewardQueue{}
		rq.Add(first)
		rq.Add(second)
		pvalequal(t, first, rq.Next())
		pvalequal(t, second, rq.Next())
		pvalequal(t, first, rq.Next())
		pvalequal(t, second, rq.Next())

	})

	t.Run("should remove address from round-robin", func(t *testing.T) {
		first := "first"
		second := "second"
		rq := RewardQueue{}
		rq.Add(first)
		rq.Add(second)
		assert.EqualValues(t, &first, rq.Next())
		rq.Remove(second)
		assert.EqualValues(t, &first, rq.Next())
	})

	// TODO test corner cases
}
