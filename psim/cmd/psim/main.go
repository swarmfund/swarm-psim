package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"

	// import services for side effects

	_ "gitlab.com/swarmfund/psim/psim/balancereporter"
	_ "gitlab.com/swarmfund/psim/psim/eventsubmitter"
	_ "gitlab.com/swarmfund/psim/psim/pokeman"

	// derivers
	_ "gitlab.com/swarmfund/psim/psim/externalsystems/btc"
	_ "gitlab.com/swarmfund/psim/psim/externalsystems/eth"

	// deposits
	_ "gitlab.com/swarmfund/psim/psim/deposits/btcdeposit"
	_ "gitlab.com/swarmfund/psim/psim/deposits/btcdepositveri"
	_ "gitlab.com/swarmfund/psim/psim/deposits/erc20"
	_ "gitlab.com/swarmfund/psim/psim/deposits/eth"
	_ "gitlab.com/swarmfund/psim/psim/deposits/ethcontracts"
	_ "gitlab.com/swarmfund/psim/psim/ethsupervisor"

	// withdrawals
	_ "gitlab.com/swarmfund/psim/psim/withdrawals/btcwithdraw"
	_ "gitlab.com/swarmfund/psim/psim/withdrawals/btcwithdveri"
	_ "gitlab.com/swarmfund/psim/psim/withdrawals/dashwithdraw"

	_ "gitlab.com/swarmfund/psim/psim/withdrawals/ethwithdraw"
	_ "gitlab.com/swarmfund/psim/psim/withdrawals/ethwithdveri"

	// funnels
	_ "gitlab.com/swarmfund/psim/psim/funnels/btcfunnel"
	_ "gitlab.com/swarmfund/psim/psim/funnels/contractfunnel"
	_ "gitlab.com/swarmfund/psim/psim/funnels/ethfunnel"

	// other folks
	_ "gitlab.com/swarmfund/psim/psim/bearer"
	_ "gitlab.com/swarmfund/psim/psim/marketmaker"
	_ "gitlab.com/swarmfund/psim/psim/notifier"
	_ "gitlab.com/swarmfund/psim/psim/prices/pricesetter"
	_ "gitlab.com/swarmfund/psim/psim/prices/pricesetterveri"
	_ "gitlab.com/swarmfund/psim/psim/template_provider"
	_ "gitlab.com/swarmfund/psim/psim/wallet_cleaner"

	// airdrops
	_ "gitlab.com/swarmfund/psim/psim/airdrop/20airdrop"
	_ "gitlab.com/swarmfund/psim/psim/airdrop/earlybird"
	_ "gitlab.com/swarmfund/psim/psim/airdrop/kycairdrop"
	_ "gitlab.com/swarmfund/psim/psim/airdrop/mrefairdrop"
	_ "gitlab.com/swarmfund/psim/psim/airdrop/telegram"

	// kyc
	_ "gitlab.com/swarmfund/psim/psim/kyc/idmind"
	_ "gitlab.com/swarmfund/psim/psim/kyc/investready"

	_ "gitlab.com/swarmfund/psim/psim/request_monitor"
)

var (
	entry          = logan.New().WithField("service", "init")
	configFile     string
	configInstance conf.Config
	rootCmd        = &cobra.Command{
		Use: "psim",
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start service with all the whistles",
		Run: func(cmd *cobra.Command, args []string) {
			instance, err := app.New(configInstance)
			if err != nil {
				entry.WithError(err).Fatal("failed to init app instance")
			}

			instance.Run()

		},
	}
)

func main() {
	cobra.OnInitialize(func() {
		configInstance = conf.NewViperConfig(configFile)
		if err := configInstance.Init(); err != nil {
			entry.WithField("config", configFile).WithError(err).Fatal("failed to init config")
		}

	})
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "config file")
	rootCmd.AddCommand(runCmd)
	err := rootCmd.Execute()
	if err != nil {

		entry.WithError(err).Fatal("something bad happened")
	}

	// // This snippet is used to quick unmarsahlling of an XDR string.
	//env := xdr.TransactionResult{}
	//err := xdr.SafeUnmarshalBase64("AAAAAAAAAAAAAAAAAAAAAQAAAAAAAAADAAAAAAAAAAAAAAB8AAAAAAzz4Jdvviw2AsGupbfHplbP4jaVAfQz4RHtZuwu6ZbaAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", &env)
	//if err != nil {
	//	panic(err)
	//}
	//bytes, err := json.Marshal(&env)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf(string(bytes))
}
