package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/psim/psim/app"
	"gitlab.com/tokend/psim/psim/conf"

	// import services for side effects
	// supervisors
	_ "gitlab.com/tokend/psim/psim/btcsupervisor"
	_ "gitlab.com/tokend/psim/psim/ethsupervisor"
	_ "gitlab.com/tokend/psim/psim/stripesupervisor"
	// other folks
	_ "gitlab.com/tokend/psim/psim/charger"
	_ "gitlab.com/tokend/psim/psim/ratesync"
	_ "gitlab.com/tokend/psim/psim/stripeverify"
	_ "gitlab.com/tokend/psim/psim/taxman"
	_ "gitlab.com/tokend/psim/psim/notifier"
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
			//env := xdr.TransactionEnvelope{}
			//err := xdr.SafeUnmarshalBase64("AAAAAACURNmoe6m4CKulrWOG62dnlfTMG13kwO3mRadEyT+1AAAAAPUfnCsAAAAAAAAAAAAAAABaAGUkAAAAAAAAAAEAAAAAAAAAEgAAAAAAlETZqHupuAirpa1jhutnZ5X0zBtd5MDt5kWnRMk/tQAAAAD1+vsyVDRGgnE/LTzfoJs6xOSkHGPapiBFK52lRFiVtwAAAAAAAAAAAAAABgAAAAAACDpPAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAF5DQq3AAAAQNSyIQiqRP2+eDpKScfQd91Bv5v0Z+dHLoWPST2ZvVSU7+o/Geg3f+hGoNZ/5cte88vr0e9FDZzOKWktiL9srAo=", &env)
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
