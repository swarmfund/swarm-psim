// Horizon requester provides helper over horizon-connector
// to make requests to Horizon not so often.
package horizonreq

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/swarmfund/go/keypair"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"context"
)

type HorizonRequester struct {
	ticker  *time.Ticker
	horizon HorizonRequestSigner
	signer  keypair.KP
}

type HorizonRequestSigner interface {
	SignedRequest(method, endpoint string, kp keypair.KP) (*http.Request, error)
}

func NewHorizonRequester(horizon HorizonRequestSigner, signer keypair.KP) func(ctx context.Context, method, url string, target interface{}) error {
	return HorizonRequester{
		ticker:  time.NewTicker(1 * time.Second),
		horizon: horizon,
		signer:  signer,
	}.do
}

func (r HorizonRequester) do(ctx context.Context, method, url string, target interface{}) error {
	// TODO proper backoff implementation
	select {
	case <-ctx.Done():
		return nil
	case <-r.ticker.C:
		// Can continue.
	}

	request, err := r.horizon.SignedRequest(method, url, r.signer)
	if err != nil {
		return errors.Wrap(err, "failed to init request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&target); err != nil {
		return errors.Wrap(err, "Failed to unmarshal response body")
	}
	return nil
}
