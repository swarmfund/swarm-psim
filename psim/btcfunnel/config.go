package btcfunnel

import "time"

type Config struct {
	ExtendedPrivateKey string `fig:"extended_private_key"`
	KeysToDerive       uint64 `fig:"keys_to_derive"`

	HotAddress  string `fig:"hot_address"`
	ColdAddress string `fig:"cold_address"`

	LastProcessedBlock uint64  `fig:"last_processed_block"`
	MinFunnelAmount    float64 `fig:"min_funnel_amount"`
	MaxHotStock        float64 `fig:"max_hot_stock"`
	DustOutputLimit    float64 `fig:"dust_output_limit"`
	BlocksToBeIncluded uint    `fig:"blocks_to_be_included"` // From 2 to 25
	MaxFeePerKB        float64 `fig:"max_fee_per_kb"`

	MinBalanceAlarmThreshold float64       `fig:"min_balance_alarm_threshold"`
	MinBalanceAlarmPeriod    time.Duration `fig:"min_balance_alarm_period"`
}
