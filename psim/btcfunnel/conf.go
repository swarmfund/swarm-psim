package btcfunnel

type Config struct {
	ExtendedPrivateKey string `fig:"extended_private_key"`
	KeysToDerive       uint64 `fig:"keys_to_derive"`

	HotAddress  string `fig:"hot_address"`
	ColdAddress string `fig:"cold_address"`

	LastProcessedBlock       uint64  `fig:"last_processed_block"`
	MinFunnelAmount          float64 `fig:"min_funnel_amount"`
	MaxHotStock              float64 `fig:"max_hot_stock"`
	DustOutputLimit          float64 `fig:"dust_output_limit"`
	MinBalanceAlarmThreshold float64 `fig:"min_amount_alarm_threshold"`

	// TODO Somehow
	//FeePerKB        float64 `fig:"fee_per_kb"`
}
