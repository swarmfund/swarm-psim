package conf

// All services existing in PSIM
const (
	// eth
	ServiceETHSupervisor      = "eth_supervisor"
	ServiceETHFunnel          = "eth_funnel"
	ServiceETHWithdraw        = "eth_withdraw"
	ServiceETHWithdrawVerify  = "eth_withdraw_verify"
	ServiceERC20Deposit       = "erc20_deposit"
	ServiceERC20DepositVerify = "erc20_deposit_verify"

	// btc
	ServiceBTCDeposit        = "btc_deposit"
	ServiceBTCDepositVerify  = "btc_deposit_verify"
	ServiceBTCFunnel         = "btc_funnel"
	ServiceBTCWithdraw       = "btc_withdraw"
	ServiceBTCWithdrawVerify = "btc_withdraw_verify"

	// Verifies
	ServiceStripeVerify = "stripe_verify"

	ServiceOperationNotifier = "notifier"
	ServiceBearer            = "bearer"
	ServicePriceSetter       = "price_setter"
	ServicePriceSetterVerify = "price_setter_verify"

	ServiceAirdropEarlybird      = "airdrop_earlybird"
	ServiceAirdropKYC            = "airdrop_kyc"
	ServiceAirdropMarchReferrals = "airdrop_march_referrals"
	ServiceAirdropMarch20b20     = "airdrop_march_20_20"

	ServiceTemplateProvider = "template_provider"
	ServiceWalletCleaner    = "wallet_cleaner"

	ServiceIdentityMind = "identity_mind"
	ServiceInvestReady  = "invest_ready"
	ServiceMixpanel     = "mixpanel"
)

// Services returns `services` slice from config, which describes, which Services to run.
func (c *ViperConfig) Services() []string {
	return c.viper.GetStringSlice("services")
}
