package btcfunnel

import "time"

type Config struct {
	ExtendedPrivateKey string `fig:"extended_private_key,required"`
	KeysToDerive       uint64 `fig:"keys_to_derive,required"`

	HotAddress  string `fig:"hot_address,required"`
	ColdAddress string `fig:"cold_address,required"`

	LastProcessedBlock uint64  `fig:"last_processed_block,required"`
	MinFunnelAmount    float64 `fig:"min_funnel_amount,required"`
	MaxHotStock        float64 `fig:"max_hot_stock,required"`
	DustOutputLimit    float64 `fig:"dust_output_limit,required"`
	BlocksToBeIncluded uint    `fig:"blocks_to_be_included,required"` // From 2 to 25
	MaxFeePerKB        float64 `fig:"max_fee_per_kb,required"`
	OffchainBlockchain string  `fig:"offchain_blockchain,required"`

	MinBalanceAlarmThreshold float64       `fig:"min_balance_alarm_threshold,required"`
	MinBalanceAlarmPeriod    time.Duration `fig:"min_balance_alarm_period,required"`
}
