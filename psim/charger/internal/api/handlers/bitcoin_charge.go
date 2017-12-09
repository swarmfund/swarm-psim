package handlers

import (
	"gitlab.com/tokend/psim/ape"
	"gitlab.com/tokend/psim/ape/problems"
	"net/http"
	"gitlab.com/tokend/psim/psim/charger/internal/resource"
)


func BitcoinChargeHandler(w http.ResponseWriter, r *http.Request) {
	body := resource.BitcoinChargeRequest{}
	if err := ape.Bind(r, &body); err != nil {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	// TODO Retrieve BTC address by AccountID from Horizon
	//globalConfig := app.Config()

	// TODO Add retrieving Charge expiration time from config

	// TODO Save amount(!), BTC address and expiration time into some storage

	// TODO Return BTC invoice (address + amount)

	w.WriteHeader(http.StatusNotImplemented)
}
