package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"

	// import services for side effects

	// eth
	_ "gitlab.com/swarmfund/psim/psim/deposits/erc20"
	_ "gitlab.com/swarmfund/psim/psim/ethfunnel"
	_ "gitlab.com/swarmfund/psim/psim/ethsupervisor"
	//_ "gitlab.com/swarmfund/psim/psim/ethwithdraw"
	_ "gitlab.com/swarmfund/psim/psim/withdrawals/eth"

	// btc
	_ "gitlab.com/swarmfund/psim/psim/btcdeposit"
	_ "gitlab.com/swarmfund/psim/psim/btcdepositveri"
	_ "gitlab.com/swarmfund/psim/psim/btcfunnel"
	_ "gitlab.com/swarmfund/psim/psim/btcwithdraw"
	_ "gitlab.com/swarmfund/psim/psim/btcwithdveri"

	// other folks
	_ "gitlab.com/swarmfund/psim/psim/airdrop"
	_ "gitlab.com/swarmfund/psim/psim/bearer"
	_ "gitlab.com/swarmfund/psim/psim/notifier"
	_ "gitlab.com/swarmfund/psim/psim/prices/pricesetter"
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
}
