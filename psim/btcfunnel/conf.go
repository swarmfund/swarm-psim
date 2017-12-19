package btcfunnel

type Config struct {
	FunnelAddress   string  `fig:"address"`
	MinFunnelAmount float64 `fig:"min_funnel_amount"`
}
