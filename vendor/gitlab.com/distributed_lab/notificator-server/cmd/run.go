package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/server"
)

var (
	Run = &cobra.Command{
		Use:   "run",
		Short: "Start service with all the whistles",
		Run: func(cmd *cobra.Command, args []string) {
			app := server.NewApp()
			q.Init()
			app.AddService(server.NewAPIService())
			app.AddService(server.NewWorkerHerder())
			app.Start()
		},
	}
)

func init() {
	Root.AddCommand(Run)
}
