package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/server"
)

var (
	Migrate = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database schema",
		Run: func(cmd *cobra.Command, args []string) {
			server.InitConf(configFile)
			q.Migrate(migrations)
		},
	}
)

func init() {
	Root.AddCommand(Migrate)
}
