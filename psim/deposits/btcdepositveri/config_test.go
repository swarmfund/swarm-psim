package btcdepositveri

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/tokend/keypair"
)

func TestNewConfig(t *testing.T) {
	t.Run("Success to create config for btc deposit verify", func(t *testing.T) {
		configData := map[string]interface{}{
			"host":                  "localhost",
			"port":                  8102,
			"last_blocks_not_watch": 15,
			"min_deposit_amount":    500000,
			"deposit_asset":         "BTC",
			"fixed_deposit_fee":     100000,
			"signer":                "SAJMOHVPENU2JPK34VR5MXH72LZZX2TABAJLHY5RQ5CXG6XVCMZU3I3N",
		}

		expected := Config{
			Host:               "localhost",
			Port:               8102,
			DepositAsset:       "BTC",
			MinDepositAmount:   500000,
			FixedDepositFee:    100000,
			LastBlocksNotWatch: 15,
			Signer:             keypair.MustParseSeed("SAJMOHVPENU2JPK34VR5MXH72LZZX2TABAJLHY5RQ5CXG6XVCMZU3I3N"),
		}

		got, err := NewConfig(configData)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, *got)
	})
	t.Run("Failed to set config, because of not enough data", func(t *testing.T) {
		configData := map[string]interface{}{}
		_, err := NewConfig(configData)
		assert.Error(t, err)
	})
}
