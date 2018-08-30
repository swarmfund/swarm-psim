package conf

// All services existing in PSIM
const (
	// integration tests for user-facing flows
	PokemanETHDepositService = "pokeman_eth_deposit"

	BalanceReporterService = "balance_reporter"
	EventSubmitterService  = "event_submitter"

	// external deployers
	ServiceBTCDeployer = "external_btc_deployer"
	ServiceETHDeployer = "external_eth_deployer"

	// deposits
	ServiceBTCDeposit       = "btc_deposit"
	ServiceBTCDepositVerify = "btc_deposit_verify"
	ServiceETHDeposit       = "eth_deposit"
	ServiceETHDepositVerify = "eth_deposit_verify"
	// DEPRECATED: eth_deposit/eth_deposit_verify is a proper implementation, left just in case
	ServiceETHSupervisor      = "eth_supervisor"
	ServiceERC20Deposit       = "erc20_deposit"
	ServiceERC20DepositVerify = "erc20_deposit_verify"
	ServiceETHContracts       = "eth_contracts_deploy"

	// funnels
	ServiceETHFunnel         = "eth_funnel"
	ServiceBTCFunnel         = "btc_funnel"
	ServiceETHContractFunnel = "eth_contract_funnel"

	// withdrawals
	ServiceBTCWithdraw       = "btc_withdraw"
	ServiceBTCWithdrawVerify = "btc_withdraw_verify"
	ServiceDashWithdraw      = "dash_withdraw"

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
	ServiceAirdropMasternode     = "airdrop_masternode"
	ServiceAirdropEarlybird      = "airdrop_earlybird"
	ServiceAirdropKYC            = "airdrop_kyc"
	ServiceAirdropMarchReferrals = "airdrop_march_referrals"
	ServiceAirdropMarch20b20     = "airdrop_march_20_20"
	ServiceAirdropTelegram       = "airdrop_telegram"

	// kyc
	ServiceIdentityMind = "identity_mind"
	ServiceInvestReady  = "invest_ready"

	ServiceTemplateProvider = "template_provider"
	ServiceWalletCleaner    = "wallet_cleaner"

	ServiceRequestMonitor = "request_monitor"
	ServiceMarketMaker    = "market_maker"
)

// Services returns `services` slice from config, which describes, which Services to run.
func (c *ViperConfig) Services() []string {
	return c.viper.GetStringSlice("services")
}
