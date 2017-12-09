package conf

// All services existing in PSIM
const (
	// Supervisors
	ServiceStripeSupervisor = "stripe_supervisor"
	ServiceBTCSupervisor    = "bitcoin_supervisor"
	ServiceETHSupervisor    = "eth_supervisor"

	// Verifies
	ServiceStripeVerify = "stripe_verify"
	ServiceBTCVerify    = "btc_verify"

	ServiceRateSync          = "rate_sync"
	ServiceTaxman            = "taxman"
	ServiceCharger           = "charger"
	ServiceOperationNotifier = "notifier"
)

// Services returns `services` slice from config, which describes, which Services to run.
func (c *ViperConfig) Services() []string {
	return c.viper.GetStringSlice("services")
}
