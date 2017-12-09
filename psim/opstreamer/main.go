package create_account_streamer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/psim/psim/app"
)

const (
	runnerName      = "operations_streamer"
	operationsLimit = 100
)

// HorizonRequestSigner is interface which is required to parametrize Streamer.
// HorizonRequestSigner is usually implemented by HorizonConnector.
type HorizonRequestSigner interface {
	SignedRequest(method, endpoint string, kp keypair.KP) (*http.Request, error)
}

// Streamer obtains Operations from time to time and streams
// found Operations into channel.
// To make Streamer start streaming, call the blocking Run() method.
type Streamer struct {
	log           *logan.Entry
	runPeriod     time.Duration
	operationType OperationType

	horizon HorizonRequestSigner
	signer  keypair.KP

	next   string
	stream chan map[string]interface{}
}

// New is constructor for Streamer
func New(log *logan.Entry, operationType OperationType, signer keypair.KP, horizon HorizonRequestSigner, runPeriod time.Duration) *Streamer {
	s := &Streamer{
		log:           log.WithField("runner", runnerName),
		runPeriod:     runPeriod,
		operationType: operationType,

		horizon: horizon,
		signer:  signer,

		next:   fmt.Sprintf("/operations?order=asc&limit=%d&operation_type=%d", operationsLimit, operationType),
		stream: make(chan map[string]interface{}),
	}

	return s
}

// GetOperationsStream returns channel where obtained Operations will be streamed into.
// Operations are streamed as key-value map of unmarshalled Operation.
// This channel is never been closed or nil.
func (s *Streamer) GetOperationsStream() <-chan map[string]interface{} {
	return s.stream
}

// Run is blocking method.
func (s *Streamer) Run(ctx context.Context) {
	app.RunOverIncrementalTimer(ctx, s.log, runnerName, s.fetchOnce, s.runPeriod)
}

func (s *Streamer) fetchOnce() error {
	next, err := s.fetchPage()
	if err != nil {
		return errors.Wrap(err, "Failed to fetch single page")
	}

	nextURL, err := url.Parse(next)
	if err != nil {
		return errors.Wrap(err, "Failed to parse URL from next", logan.Field("next", next))
	}

	if nextURL.RequestURI() == s.next {
		// Finished - this page is the last one.
		return nil
	}

	s.next = nextURL.RequestURI()
	return nil
}

func (s *Streamer) fetchPage() (next string, err error) {
	request, err := s.horizon.SignedRequest("GET", s.next, s.signer)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create signed request")
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "Failed to do request")
	}

	defer response.Body.Close()

	var body OperationsResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal OperationsResponse from response from Horizon")
	}

	for _, op := range body.Embedded.Records {
		s.stream <- op
	}

	s.log.WithField("count", len(body.Embedded.Records)).Debug("Fetched operations page from Horizon.")

	return body.Links.Next.HREF, nil
}
