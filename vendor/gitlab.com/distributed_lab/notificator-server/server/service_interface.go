package server

import (
	"gitlab.com/distributed_lab/notificator-server/conf"
)

type ServiceInterface interface {
	Init(cfg conf.Config)
	Run(cfg conf.Config)
}
