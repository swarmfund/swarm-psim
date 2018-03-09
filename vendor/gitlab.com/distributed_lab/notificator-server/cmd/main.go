package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/notificator-server/auth"
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/q"
	"gitlab.com/distributed_lab/notificator-server/server"
)

var (
	configFile     string
	configInstance conf.Config
	migrations     string
	entry          = *logrus.New()
	rootCmd        = &cobra.Command{
		Use: "notificator",
	}
	createMigrationCmd = &cobra.Command{
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
			if err := q.NewMigration(configInstance.DB().DSN, migrations, name); err != nil {
				entry.WithError(err).Fatal("create migration failed")
			}
		},
	}
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate database schema",
		Run: func(cmd *cobra.Command, args []string) {
			q.Migrate(configInstance.DB().DSN, migrations, configInstance.Log())
		},
	}
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Start service with all the whistles",
		Run: func(cmd *cobra.Command, args []string) {
			app := server.NewApp()
			q.Init(configInstance.DB().Driver, configInstance.DB().DSN, configInstance.Log())
			app.AddService(server.NewAPIService(configInstance.Log(), configInstance.Requests()))
			app.AddService(server.NewWorkerHerder(configInstance))
			app.Start(configInstance)
		},
	}
	tokensCmd = &cobra.Command{
		Use:   "tokens",
		Short: "Print new client token pair",
		Run: func(cmd *cobra.Command, args []string) {
			q.Init(configInstance.DB().Driver, configInstance.DB().DSN, configInstance.Log())
			pair, err := auth.GeneratePair()
			if err != nil {
				entry.WithError(err).Fatal("failed to generate pair")
			}

			if err = q.NewQ(configInstance.DB().Driver, configInstance.DB().DSN, configInstance.Log()).Auth().Insert(pair); err != nil {
				entry.WithError(err).Fatal("failed to save pair")
			}
			fmt.Println("public: ", pair.Public)
			fmt.Println("secret: ", pair.Secret)
		},
	}
)

func initConfig(fn string) (conf.Config, error) {
	c := conf.NewViperConfig(fn)
	if err := c.Init(); err != nil {
		return nil, err
	}
	return c, nil
}

func main() {
	cobra.OnInitialize(func() {
		c, err := initConfig(configFile)
		if err != nil {
			entry.WithField("service", "init").WithError(err).Fatal("failed to init config")
		}
		configInstance = c

	})

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVar(&migrations, "migrations", "", "Migrations dir")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(createMigrationCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(tokensCmd)

	if err := rootCmd.Execute(); err != nil {
		entry.WithField("service", "notificator").WithError(err).Fatal("something bad happened")
	}
}
