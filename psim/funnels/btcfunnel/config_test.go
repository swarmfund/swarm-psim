package btcfunnel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/swarmfund/psim/psim/supervisor"
)

func TestBTCFunnelConfig(t *testing.T) {
	configData := map[string]interface{}{
		"extended_private_key":        "xprv...",
		"keys_to_derive":              "10000",
		"hot_address":                 "2N1w4RzejEWkCyumsZY8prvmRxAPFkcwehb",
		"cold_address":                "2N8hwP1WmJrFF5QWABn38y63uYLhnJYJYTF",
		"last_processed_block":        "1260685",
		"min_funnel_amount":           "0",
		"max_hot_stock":               "50",
		"dust_output_limit":           "0.005",
		"max_fee_per_kb":              "0.00024",
		"blocks_to_be_included":       "4",
		"network_type":                "1",
		"min_balance_alarm_threshold": "100",
		"min_balance_alarm_period":    "1h",
		"disable_low_balance_monitor": "true",
	}

	config := Config{}

	err := figure.
		Out(&config).
		From(configData).
		With(figure.BaseHooks, supervisor.DLFigureHooks).
		Please()
	assert.NoError(t, err)
	assert.True(t, config.DisableLowBalanceMonitor)
}
