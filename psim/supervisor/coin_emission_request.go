package supervisor

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
)

type IssuanceRequestOpt struct {
	Reference string
	Receiver  string
	Asset     string
	Amount    uint64
	Details   string
}

func (s *Service) CraftIssuanceRequest(opt IssuanceRequestOpt) *xdrbuild.Transaction {
	return s.builder.
		Transaction(s.config.ExchangeKP).
		Op(xdrbuild.CreateIssuanceRequestOp{
			Reference: opt.Reference,
			Receiver:  opt.Receiver,
			Asset:     opt.Asset,
			Amount:    opt.Amount,
			Details:   opt.Details,
		}).
		Sign(s.config.SignerKP)
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
		return errors.From(errors.New("Got unsuccessful response from neighbour-verifier"), logan.F{"status_code": response.StatusCode})
	}
	return nil
}
