package ethcontracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/keypair"
)

func TestNewConfig(t *testing.T) {
	kp, _ := keypair.Random()
	valid := map[string]interface{}{
		"source":         kp.Address(),
		"signer":         kp.Seed(),
		"target_count":   10,
		"external_types": []int32{4, 2},
	}

	t.Run("valid", func(t *testing.T) {
		expected := Config{
			TargetCount: 10,
			Source:      keypair.MustParseAddress(kp.Address()),
			Signer:      kp,
			//ExternalTypes: []int32{4, 2},
		}
		got, err := NewConfig(valid)
		assert.NoError(t, err)
		assert.Equal(t, &expected, got)
	})
}
