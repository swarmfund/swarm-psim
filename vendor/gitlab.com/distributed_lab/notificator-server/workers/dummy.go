package workers

import (
	"gitlab.com/distributed_lab/notificator-server/conf"
	"gitlab.com/distributed_lab/notificator-server/types"
)

func Dummy(_ types.Request, _ conf.Config) bool {
	return true
}
