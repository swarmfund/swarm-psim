package server

import "gitlab.com/distributed_lab/notificator-server/types"

type TaskResult struct {
	Request types.Request
	Success bool
}
