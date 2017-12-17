package conf

// All services existing in PSIM
const (
	// eth
	ServiceETHSupervisor = "eth_supervisor"
	ServiceETHFunnel     = "eth_funnel"

	// btc
	ServiceBTCSupervisor = "bitcoin_supervisor"

	// Verifies
	ServiceStripeVerify = "stripe_verify"
	ServiceBTCVerify    = "btc_verify"

	ServiceRateSync          = "rate_sync"
	ServiceCharger           = "charger"
	ServiceOperationNotifier = "notifier"
)

// Services returns `services` slice from config, which describes, which Services to run.
func (c *ViperConfig) Services() []string {
	return c.viper.GetStringSlice("services")
}
