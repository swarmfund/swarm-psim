package main

import (
	commands "gitlab.com/distributed_lab/notificator-server/cmd"
	"gitlab.com/distributed_lab/notificator-server/log"
)

func main() {
	if err := commands.Root.Execute(); err != nil {
		log.WithField("service", "notificator").WithError(err).Fatal("something bad happened")
	}
}
