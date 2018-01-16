package conf

// All services existing in PSIM
const (
	// eth
	ServiceETHSupervisor = "eth_supervisor"
	ServiceETHFunnel     = "eth_funnel"
	ServiceETHWithdraw   = "eth_withdraw"

	// btc
	ServiceBTCSupervisor     = "bitcoin_supervisor"
	ServiceBTCVerify         = "btc_verify"
	ServiceBTCFunnel         = "btc_funnel"
	ServiceBTCWithdraw       = "btc_withdraw"
	ServiceBTCWithdrawVerify = "btc_withdraw_verify"

	// Verifies
	ServiceStripeVerify = "stripe_verify"

	ServiceCharger           = "charger"
	ServiceOperationNotifier = "notifier"
	ServiceBearer            = "bearer"
)

// Services returns `services` slice from config, which describes, which Services to run.
func (c *ViperConfig) Services() []string {
	return c.viper.GetStringSlice("services")
}
