package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/notificator-server/q"
)

var (
	CreateMigration = &cobra.Command{
		Use:   "createmigration",
		Short: "Creates new migrations files",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				entry.Fatal("too many arguments")
			}
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			if err := q.NewMigration(migrations, name); err != nil {
				entry.WithError(err).Fatal("create migration failed")
			}
		},
	}
)

func init() {
	Root.AddCommand(CreateMigration)
}
