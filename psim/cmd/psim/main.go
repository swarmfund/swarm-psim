package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"

	// import services for side effects
	// eth
	_ "gitlab.com/swarmfund/psim/psim/ethfunnel"
	_ "gitlab.com/swarmfund/psim/psim/ethsupervisor"
	// btc
	//_ "gitlab.com/swarmfund/psim/psim/btcsupervisor"
	// other folks
	_ "gitlab.com/swarmfund/psim/psim/charger"
	_ "gitlab.com/swarmfund/psim/psim/notifier"
	_ "gitlab.com/swarmfund/psim/psim/ratesync"
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
			//err := xdr.SafeUnmarshalBase64("AAAAAP4DpOIcoI8urCJITRZtEDS0wzyPuGojb7AbKpHcMR1gAAAAAAAAAAAAAAAAAAAAAAAAAABaOnMWAAAAAAAAAAEAAAAAAAAAAwAAAANTVU4AAAAAAAAAAAAAAAAAZdVETc5jpkVXYmmpABfBQzkkERKbmqpjtJAPalBBIpcAAAAAAAAAQjB4ODZmOGQwZmZlNmI1MDI5MDJhZmFhYTJiZDA2ODYzYWEwM2JjZDhkZjBiNWI3NmJkYzMxYjM5OGNmZTBjY2QwMQAAAAAAAAAAAAAAAAAB3DEdYAAAAEA6mRTtWYjxpCcim66actqKdwdqGbUt7N+VPxiaBbtk8TCYHanMIHSl9esYkTrwS4qU4gWhj4kZbK3D6azWgZcJ", &env)
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
