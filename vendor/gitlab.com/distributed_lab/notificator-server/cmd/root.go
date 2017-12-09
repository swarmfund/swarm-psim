package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/notificator-server/log"
	"gitlab.com/distributed_lab/notificator-server/server"
)

var (
	configFile string
	migrations string
	entry      = log.WithField("service", "notificator")
	Root       = &cobra.Command{
		Use: "notificator",
	}
)

func init() {
	cobra.OnInitialize(func() {
		server.InitConf(configFile)
	})
	Root.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	Root.PersistentFlags().StringVar(&migrations, "migrations", "", "Migrations dir")

	cobra.OnInitialize(func() {
		server.InitConf(configFile)
	})
}
