package btcfunnel

import (
	"time"

	"gitlab.com/swarmfund/psim/psim/externalsystems/derive"
)

type Config struct {
	ExtendedPrivateKey       string             `fig:"extended_private_key,required"`
	KeysToDerive             uint64             `fig:"keys_to_derive,required"`
	HotAddress               string             `fig:"hot_address,required"`
	ColdAddress              string             `fig:"cold_address,required"`
	LastProcessedBlock       uint64             `fig:"last_processed_block,required"`
	MinFunnelAmount          float64            `fig:"min_funnel_amount,required"`
	MaxHotStock              float64            `fig:"max_hot_stock,required"`
	DustOutputLimit          float64            `fig:"dust_output_limit,required"`
	BlocksToBeIncluded       uint               `fig:"blocks_to_be_included,required"` // From 2 to 25
	MaxFeePerKB              float64            `fig:"max_fee_per_kb,required"`
	NetworkType              derive.NetworkType `fig:"network_type,required"`
	MinBalanceAlarmThreshold float64            `fig:"min_balance_alarm_threshold,required"`
	MinBalanceAlarmPeriod    time.Duration      `fig:"min_balance_alarm_period,required"`
}

func (c Config) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"keys_to_derive":              c.KeysToDerive,
		"hot_address":                 c.HotAddress,
		"cold_address":                c.ColdAddress,
		"last_processed_block":        c.LastProcessedBlock,
		"min_funnel_amount":           c.MinFunnelAmount,
		"max_hot_stock":               c.MaxHotStock,
		"dust_output_limit":           c.DustOutputLimit,
		"blocks_to_be_included":       c.BlocksToBeIncluded,
		"max_fee_per_kb":              c.MaxFeePerKB,
		"network_type":                c.NetworkType,
		"min_balance_alarm_threshold": c.MinBalanceAlarmThreshold,
		"min_balance_alarm_period":    c.MinBalanceAlarmPeriod,
	}
}
