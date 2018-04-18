package conf

// All services existing in PSIM
const (
	// deposits
	ServiceBTCDeposit         = "btc_deposit"
	ServiceBTCDepositVerify   = "btc_deposit_verify"
	ServiceETHSupervisor      = "eth_supervisor"
	ServiceERC20Deposit       = "erc20_deposit"
	ServiceERC20DepositVerify = "erc20_deposit_verify"
	ServiceETHContracts       = "eth_contracts_deploy"

	// funnels
	ServiceETHFunnel = "eth_funnel"
	ServiceBTCFunnel = "btc_funnel"

	// withdrawals
	ServiceBTCWithdraw       = "btc_withdraw"
	ServiceBTCWithdrawVerify = "btc_withdraw_verify"
	ServiceETHWithdraw       = "eth_withdraw"
	ServiceETHWithdrawVerify = "eth_withdraw_verify"

	// Verifies
	ServiceStripeVerify = "stripe_verify"

	ServiceOperationNotifier = "notifier"
	ServiceBearer            = "bearer"

	// prices
	ServicePriceSetter       = "price_setter"
	ServicePriceSetterVerify = "price_setter_verify"

	// airdrops
	ServiceAirdropEarlybird      = "airdrop_earlybird"
	ServiceAirdropKYC            = "airdrop_kyc"
	ServiceAirdropMarchReferrals = "airdrop_march_referrals"
	ServiceAirdropMarch20b20     = "airdrop_march_20_20"

	// kyc
	ServiceIdentityMind = "identity_mind"

	ServiceTemplateProvider = "template_provider"
	ServiceWalletCleaner    = "wallet_cleaner"

	ServiceIdentityMind = "identity_mind"
	ServiceInvestReady  = "invest_ready"
)

// Services returns `services` slice from config, which describes, which Services to run.
func (c *ViperConfig) Services() []string {
	return c.viper.GetStringSlice("services")
}
