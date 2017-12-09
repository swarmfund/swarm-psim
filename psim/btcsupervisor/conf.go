package btcsupervisor

import "gitlab.com/swarmfund/psim/psim/supervisor"

type Config struct {
	Supervisor supervisor.Config `fig:"supervisor"`

	LastProcessedBlock uint64 `fig:"last_processed_block"`
}
