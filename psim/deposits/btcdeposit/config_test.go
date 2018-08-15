package btcdeposit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/psim/psim/supervisor"
	"gitlab.com/tokend/keypair"
)

func TestNewConfig(t *testing.T) {
	t.Run("Successfully set config", func(t *testing.T) {
		configData := map[string]interface{}{
			"source":                "GB4GMHDGROECF4J4XJ7TJUVMGT4IHVFT2245QDOTNQSLIDL2HAMSMJJA",
			"signer":                "SCDHB3LSE24SKIWVT3DRMN3RATR53TOCAOENH5SK6UW2M7EGOPGK2KFX",
			"fixed_deposit_fee":     100000,
			"last_blocks_not_watch": 15,
			"deposit_asset":         "BTC",
			"last_processed_block":  1260685,
			"min_deposit_amount":    500000,
		}

		expected := Config{
			Supervisor:         supervisor.Config{},
			LastProcessedBlock: 1260685,
			MinDepositAmount:   500000,
			DepositAsset:       "BTC",
			FixedDepositFee:    100000,
			Signer:             keypair.MustParseSeed("SCDHB3LSE24SKIWVT3DRMN3RATR53TOCAOENH5SK6UW2M7EGOPGK2KFX"),
			Source:             keypair.MustParseAddress("GB4GMHDGROECF4J4XJ7TJUVMGT4IHVFT2245QDOTNQSLIDL2HAMSMJJA"),
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
