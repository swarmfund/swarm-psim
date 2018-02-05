package conf

// All services existing in PSIM
const (
	// eth
	ServiceETHSupervisor = "eth_supervisor"
	ServiceETHFunnel     = "eth_funnel"
	ServiceETHWithdraw   = "eth_withdraw"

	// btc

	// DEPRECATED use ServiceBTCDeposit instead
	ServiceBTCSupervisor     = "bitcoin_supervisor"
	ServiceBTCDeposit        = "btc_deposit"

	ServiceBTCVerify         = "btc_verify"
	ServiceBTCFunnel         = "btc_funnel"
	ServiceBTCWithdraw       = "btc_withdraw"
	ServiceBTCWithdrawVerify = "btc_withdraw_verify"

	// Verifies
	ServiceStripeVerify = "stripe_verify"

	ServiceOperationNotifier = "notifier"
	ServiceBearer            = "bearer"
)

// Services returns `services` slice from config, which describes, which Services to run.
func (c *ViperConfig) Services() []string {
	return c.viper.GetStringSlice("services")
}
