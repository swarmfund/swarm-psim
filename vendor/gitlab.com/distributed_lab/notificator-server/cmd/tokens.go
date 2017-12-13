package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/distributed_lab/notificator-server/auth"
	"gitlab.com/distributed_lab/notificator-server/q"
)

var (
	Tokens = &cobra.Command{
		Use:   "tokens",
		Short: "Print new client token pair",
		Run: func(cmd *cobra.Command, args []string) {
			q.Init()
			pair, err := auth.GeneratePair()
			if err != nil {
				entry.WithError(err).Fatal("failed to generate pair")
			}
			if err = q.NewQ().Auth().Insert(pair); err != nil {
				entry.WithError(err).Fatal("failed to save pair")
			}
			fmt.Println("public: ", pair.Public)
			fmt.Println("secret: ", pair.Secret)
		},
	}
)

func init() {
	Root.AddCommand(Tokens)
}
