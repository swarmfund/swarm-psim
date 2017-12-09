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
	runnerName = "create_account_streamer"
)

var (
	justClosedChan = make(chan struct{})
)

func init() {
	close(justClosedChan)
}

type HorizonRequestSigner interface {
	SignedRequest(method, endpoint string, kp keypair.KP) (*http.Request, error)
}

type CreatedAccountsStreamer struct {
	log       *logan.Entry
	runPeriod time.Duration

	horizon HorizonRequestSigner
	signer  keypair.KP

	next string
	// TODO - Fix raise condition!
	errorWaiter chan struct{}
	// TODO - Fix raise condition!
	initialReady chan struct{}
	stream       chan CreateAccountOp
}

// New is constructor for CreatedAccountsStreamer
func New(log *logan.Entry, signer keypair.KP, horizon HorizonRequestSigner, runPeriod time.Duration) *CreatedAccountsStreamer {
	s := &CreatedAccountsStreamer{
		log:       log.WithField("runner", runnerName),
		runPeriod: runPeriod,

		horizon: horizon,
		signer:  signer,

		next:         fmt.Sprintf("/operations?order=asc&limit=100&operation_type=0"),
		initialReady: make(chan struct{}),
		stream:       make(chan CreateAccountOp),
	}

	return s
}

// ReadinessWaiter returns channel, which shows that streamer is ready,
// if reading from it doesn't block i.e. channel is closed.
// Will never return nil channel.
func (s *CreatedAccountsStreamer) ReadinessWaiter() <-chan struct{} {
	if s.initialReady != nil {
		return s.initialReady
	}

	// TODO - Fix raise condition!
	if s.errorWaiter != nil {
		return s.errorWaiter
	}

	return justClosedChan
}

// GetStream returns channel where obtained CreatedAccountOps will be streamed into.
// This channel is never been closed or nil.
func (s *CreatedAccountsStreamer) GetStream() <-chan CreateAccountOp {
	return s.stream
}

// Run is blocking method.
func (s *CreatedAccountsStreamer) Run(ctx context.Context) {
	app.RunOverIncrementalTimer(ctx, s.log, runnerName, s.fetchOnceSetErrorState, s.runPeriod)
}

func (s *CreatedAccountsStreamer) fetchOnceSetErrorState() error {
	err := s.fetchOnce()

	//s.errorState = err != nil
	if err != nil {
		s.errorWaiter = make(chan struct{})
		return err
	}

	// err == nil
	if s.errorWaiter != nil {
		close(s.errorWaiter)
		s.errorWaiter = nil
	}

	return nil
}

func (s *CreatedAccountsStreamer) fetchOnce() error {
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

		if s.initialReady != nil {
			close(s.initialReady)
			s.initialReady = nil
		}

		return nil
	}

	s.next = nextURL.RequestURI()
	return nil
}

func (s *CreatedAccountsStreamer) fetchPage() (next string, err error) {
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
