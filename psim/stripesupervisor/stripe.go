package stripesupervisor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stripe/stripe-go"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/stripeverify"
)

// TODO Listen to context and close Errors
func (s *Service) processStripeHistory() {
	listParams := stripe.ChargeListParams{}
	charges := s.stripe.Charges.List(&listParams)
	ticker := time.NewTicker(5 * time.Second)

	// getting earliest transaction to get cursors work way we want
	cursor := ""
	for charges.Next() {
		if app.IsCanceled(s.Ctx) {
			return
		}

		if !s.IsLeader {
			select {
			case <-ticker.C:
				continue
			case <-s.Ctx.Done():
				return
			}
		}

		// I am leader!

		charge := charges.Charge()
		cursor = charge.ID
	}

	for {
		if !s.IsLeader {
			select {
			case <-ticker.C:
				continue
			case <-s.Ctx.Done():
				return
			}
		}

		listParams = stripe.ChargeListParams{}
		listParams.End = cursor
		charges := s.stripe.Charges.List(&listParams)

		for charges.Next() {
			if app.IsCanceled(s.Ctx) {
				return
			}

			charge := charges.Charge()

			possibleCursor, err := s.processStripeCharge(charge)
			if err != nil {
				s.Errors <- errors.Wrap(err, "Failed to process charge")
				continue
			}

			if possibleCursor != nil {
				cursor = *possibleCursor
			}
		}
	}
}

// TODO split to several methods
func (s *Service) processStripeCharge(charge *stripe.Charge) (*string, error) {
	log := s.Log.WithField("charge_id", charge.ID)

	reference := charge.Meta["reference"]
	if reference == "" {
		// charge w/o reference, no-op
		return &charge.ID, nil
	}

	receiver := charge.Meta["receiver"]
	if receiver == "" {
		// charge w/o receiver, no-op
		return &charge.ID, nil
	}

	asset := charge.Meta["asset"]
	if asset == "" {
		// charge w/o asset, no-op
		return &charge.ID, nil
	}

	switch charge.Status {
	case "succeeded":
		log.Info("processing charge")
	case "failed":
		return &charge.ID, nil
	case "pending":
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown charge status: %s", charge.Status)
	}

	alreadyCounted, err := s.CheckCoinEmissionRequestExistence(reference)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to check existence of CoinEmissionRequest")
	}

	if alreadyCounted {
		return &charge.ID, nil
	}

	amount := charge.Amount * 100
	envelope, err := s.PrepareCEREnvelope(reference, receiver, asset, amount)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to craft CoinEmissionRequests envelope")
	}

	verifyPayload, err := json.Marshal(stripeverify.VerifyRequest{
		Envelope: *envelope,
		ChargeID: charge.ID,
	})

	s.SendCoinEmissionRequest(s.Ctx, verifyPayload)
	if app.IsCanceled(s.Ctx) {
		return nil, nil
	}

	return &charge.ID, nil
}
