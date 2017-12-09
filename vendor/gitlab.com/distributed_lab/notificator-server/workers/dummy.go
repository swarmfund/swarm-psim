package workers

import "gitlab.com/distributed_lab/notificator-server/types"

func Dummy(_ types.Request) bool {
	return true
}
