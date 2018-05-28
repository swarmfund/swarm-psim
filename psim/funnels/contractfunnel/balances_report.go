package contractfunnel

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

func (s *Service) printBalancesReport() {
	report, err := s.prepareBalancesReport()
	if err != nil {
		s.log.WithError(err).Error("Failed to prepare Contracts balances report.")
		return
	}

	bb, err := json.Marshal(report)
	if err != nil {
		panic(errors.Wrap(err, "Failed to marshal prepared Contracts balances report"))
		return
	}

	s.log.WithField("report", string(bb)).Info("Contracts balances report prepared.")
}

func (s *Service) prepareBalancesReport() (map[string]uint64, error) {
	report := make(map[string]uint64)

	for addr, _ := range s.contracts {
		contractBalance, err := s.erc20Contract.BalanceOf(nil, addr)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get balance of the Contract", logan.F{
				"contract_address": addr,
			})
		}

		report[addr.String()] = contractBalance.Uint64()
	}

	return report, nil
}
