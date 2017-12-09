package ethsupervisor

import (
	"gitlab.com/tokend/psim/psim/supervisor"
)

type Config struct {
	Supervisor supervisor.Config `fig:"supervisor"`
}
