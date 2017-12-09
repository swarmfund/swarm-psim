package handlers

import (
	"net/http"

	"github.com/stripe/stripe-go"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/charger/internal/resource"
)

func StripeChargeHandler(w http.ResponseWriter, r *http.Request) {
	body := resource.StripeChargeRequest{}
	if err := ape.Bind(r, &body); err != nil {
		ape.RenderErr(w, r, problems.BadRequest("Cannot decode JSON request."))
		return
	}

	chargeParams := &stripe.ChargeParams{
		Amount:   uint64(body.Amount / 100),
		Currency: "USD",
		Source: &stripe.SourceParams{
			Token: body.Token,
		},
		Params: stripe.Params{
			IdempotencyKey: body.Reference,
			Meta: map[string]string{
				"receiver":  body.Receiver,
				"reference": body.Reference,
				"asset":     body.Asset,
			},
		},
	}

	// send Charge request to Stripe
	_, err := Stripe(r).Charges.New(chargeParams)
	if err != nil {
		Log(r).WithError(err).Error("failed to send charge request to Stripe")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
}
