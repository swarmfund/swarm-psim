package ethsupervisor

import (
	"gitlab.com/swarmfund/psim/psim/supervisor"
)

type Config struct {
	Supervisor supervisor.Config `fig:"supervisor"`
}
