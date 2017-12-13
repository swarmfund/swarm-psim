package internal

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/addrstate"
)

type HorizonRequester struct {
	ticker  *time.Ticker
	horizon *horizon.Connector
	signer  keypair.KP
}

func NewHorizonRequester(horizon *horizon.Connector, signer keypair.KP) addrstate.Requester {
	return HorizonRequester{
		ticker:  time.NewTicker(5),
		horizon: horizon,
		signer:  signer,
	}.do
}

func (r HorizonRequester) do(method, endpoint string, target interface{}) error {
	// TODO proper backoff implementation
	<-r.ticker.C
	request, err := r.horizon.SignedRequest(method, endpoint, r.signer)
	if err != nil {
		return errors.Wrap(err, "failed to init request")
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&target); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}
