package pricesetterveri

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/swarmfund/psim/psim/prices/providers"
	"gitlab.com/tokend/keypair"
)

func TestNewConfig(t *testing.T) {
	t.Run("Success to create config for price setter verify", func(t *testing.T) {

		configData := map[string]interface{}{
			"signer":      "SAV2QWTC4OH44U6UF7THPJGMY2G3IE4PVSY5WM7AHEE3F2PQSNS6OEH5",
			"host":        "localhost",
			"port":        8501,
			"base_asset":  "ETH",
			"quote_asset": "USD",
			"providers": []interface{}{
				map[interface{}]interface{}{"name": "bitfinex", "period": "15s"},
			},
			"providers_to_agree":      3,
			"max_price_delta_percent": "10",
			"verifier_service_name":   "price_setter_verify",
		}

		expected := Config{
			Host:                 "localhost",
			Port:                 8501,
			BaseAsset:            "ETH",
			QuoteAsset:           "USD",
			Providers:            []providers.ProviderConfig{{Name: "bitfinex", Period: time.Second * 15}},
			ProvidersToAgree:     3,
			MaxPriceDeltaPercent: "10",
			Signer:               keypair.MustParseSeed("SAV2QWTC4OH44U6UF7THPJGMY2G3IE4PVSY5WM7AHEE3F2PQSNS6OEH5"),
			VerifierServiceName:  "price_setter_verify",
		}

		got, err := NewConfig(configData)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, *got)
	})

	t.Run("Failed to create config price setter, because of not enough data", func(t *testing.T) {
		configData := map[string]interface{}{}
		_, err := NewConfig(configData)
		assert.Error(t, err)
	})
}
