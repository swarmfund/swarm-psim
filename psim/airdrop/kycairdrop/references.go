package kycairdrop

import (
	"context"

	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/airdrop"
)

// FetchAllReferences is a blocking method
func (s *Service) fetchAllReferences(ctx context.Context) {
	running.UntilSuccess(ctx, s.log, "all_references_fetcher", func(ctx context.Context) (bool, error) {
		allReferences, err := s.referencesProvider.References(s.config.Source.Address())
		if err != nil {
			return false, errors.Wrap(err, "Failed to get References from Horizon")
		}

		for _, r := range allReferences {
			s.existingReferences = append(s.existingReferences, r.Reference)
		}

		return true, nil
	}, 10*time.Second, time.Hour)

}

func (s *Service) isAlreadyIssued(accAddress string) bool {
	reference := airdrop.BuildReference(accAddress, airdrop.KYCReferenceSuffix)

	for _, r := range s.existingReferences {
		if r == reference {
			return true
		}
	}

	return false
}
