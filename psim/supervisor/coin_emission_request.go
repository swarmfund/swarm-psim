package supervisor

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

const (
	baseAsset = "SUN"
)

// PrepareCERTx uses HorizonConnector - prepares Transaction with a CoinEmissionRequest operation,
// signs it using SignerKP and returns marshaled Transaction.
//
// Specific supervisors use this method to create Transaction to be sent to a verifier.
func (s *Service) PrepareCERTx(reference, receiver string, amount uint64) *horizon.TransactionBuilder {
	return s.horizon.Transaction(&horizon.TransactionBuilder{Source: s.config.ExchangeKP}).
		Op(&horizon.CreateIssuanceRequestOp{
			Reference: reference,
			Receiver:  receiver,
			Asset:     baseAsset,
			Amount:    amount,
		}).
		Sign(s.config.SignerKP)
}

// CheckCoinEmissionRequestExistence returns true if CER with such `reference` already exists.
func (s *Service) CheckCoinEmissionRequestExistence(reference string) (bool, error) {
	cers, err := s.horizon.CoinEmissionRequests(s.config.SignerKP, &horizon.CoinEmissionRequestsParams{
		Reference: reference,
		Exchange:  s.config.ExchangeKP.Address(),
	})
	if err != nil {
		return false, errors.Wrap(err, "Failed to get CoinEmissionRequests from Horizon")
	}

	if len(cers) > 0 {
		s.Log.Info("already counted")
		return true, nil
	}

	return false, nil
}

// SendCoinEmissionRequest only returns if success or ctx canceled.
func (s *Service) SendCoinEmissionRequestForVerify(ctx context.Context, verifyRequest []byte) {
	ticker := time.NewTicker(5 * time.Second)

	for ; true; <-ticker.C {
		if app.IsCanceled(ctx) {
			return
		}

		// FIXME STRIPE
		neighbors, err := s.discovery.DiscoverService(conf.ServiceStripeVerify)
		if err != nil {
			s.Log.WithField("service_to_discover", conf.ServiceStripeVerify).WithError(err).Error("Failed to discover neighbors.")
			continue
		}
		if len(neighbors) == 0 {
			continue
		}

		s.Log.WithField("count", len(neighbors)).Info("Discovered neighbors.")

		for _, neighbor := range neighbors {
			err := s.sendForVerifying(verifyRequest, neighbor)
			if err != nil {
				s.Log.WithField("neighbor_address", neighbor.Address).WithError(err).Error("Failed to verify via a neighbour.")
				continue
			}

			// Sent to verifier successfully
			return
		}
	}
}

func (s *Service) sendForVerifying(txToVerify []byte, neighbor discovery.ServiceEntry) error {
	response, err := http.Post(neighbor.Address, "application/psim+envelope", bytes.NewReader(txToVerify))
	if err != nil {
		return errors.New("failed to submit verification")
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.From(errors.New("Got unsuccessful response from neighbour-verifier"), logan.Field("status_code", response.StatusCode))
	}
	return nil
}
