package btcfunnel

type Config struct {
	HotAddress  string `fig:"hot_address"`
	ColdAddress string `fig:"cold_address"`

	MinFunnelAmount float64 `fig:"min_funnel_amount"`
	MaxHotStock     float64 `fig:"max_hot_stock"`
	DustOutputLimit float64 `fig:"dust_output_limit"`
}
