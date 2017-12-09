package server

import (
	"gitlab.com/distributed_lab/notificator-server/limiter"
	"gitlab.com/distributed_lab/notificator-server/types"
)

const (
	_ types.RequestTypeID = iota
	RequestTypeDummy
	RequestTypeUniqueDummy
	RequestTypeUniqueEmail
	RequestTypeUniqueSMS
)

type RequestType struct {
	ID       types.RequestTypeID
	Priority int
	Worker   string
	Limiters []limiter.Limiter
}
