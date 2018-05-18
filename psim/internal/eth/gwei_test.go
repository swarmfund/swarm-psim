package eth

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestFromGwei(t *testing.T) {
	cases := []struct {
		input    *big.Int
		expected *big.Int
	}{
		{input: big.NewInt(42), expected: big.NewInt(42000000000)},
	}
	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, tc.expected, FromGwei(tc.input))
		})
	}
}
