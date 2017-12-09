package snapshoter

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/psim/psim/taxman/internal/state"
)

type referenceInput struct {
	BalanceID state.BalanceID
	Ledger    int64
}

func TestReference(t *testing.T) {
	input := []referenceInput{
		{"balance_1", 12},
		{"balance_1", 13},
		{"balance_1", 14},
		{"balance_2", 12},
		{"balance_3", 13},
		{"balance_4", 13},
	}

	Convey("Determinism", t, func() {
		Convey("payoutTypeReferral", func() {
			checkDeterminism(t, payoutTypeReferral, input)
		})
		Convey("payoutTypeToken", func() {
			checkDeterminism(t, payoutTypeToken, input)
		})
	})
	Convey("Uniqueness", t, func() {
		referral := generateReferences(payoutTypeReferral, input)
		token := generateReferences(payoutTypeToken, input)
		for refKey := range referral {
			_, isInToken := token[refKey]
			So(isInToken, ShouldBeFalse)
		}
	})
}

func checkDeterminism(t *testing.T, payoutT payoutType, inputs []referenceInput) {
	first := generateReferences(payoutT, inputs)
	second := generateReferences(payoutT, inputs)
	assert.Equal(t, first, second)
}

func generateReferences(payoutT payoutType, inputs []referenceInput) map[string]bool {
	result := map[string]bool{}
	for _, input := range inputs {
		ref := reference(payoutT, input.BalanceID, input.Ledger)
		_, exists := result[ref]
		So(exists, ShouldBeFalse)
		result[ref] = true
	}

	return result
}
