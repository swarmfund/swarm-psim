package telegram

import (
	"context"

	"fmt"
	"net/http"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/issuance"
	"gitlab.com/tokend/go/doorman"
)

type IssuanceSubmitter interface {
	Submit(ctx context.Context, accountAddress, balanceID string, amount uint64, opDetails string) (*issuance.RequestOpt, bool, error)
}

type BalanceIDProvider interface {
	GetBalanceID(accAddress, asset string) (*string, error)
}

type Service struct {
	log    *logan.Entry
	config Config

	issuanceSubmitter IssuanceSubmitter
	balanceIDProvider BalanceIDProvider
	doorman           doorman.Doorman

	blackList map[string]struct{}
}

func NewService(
	log *logan.Entry,
	config Config,
	issuanceSubmitter IssuanceSubmitter,
	balanceIDProvider BalanceIDProvider,
	doorman doorman.Doorman) *Service {

	return &Service{
		log:    log,
		config: config,

		issuanceSubmitter: issuanceSubmitter,
		balanceIDProvider: balanceIDProvider,
		doorman:           doorman,

		blackList: make(map[string]struct{}),
	}
}

func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	for _, accID := range s.config.BlackList {
		s.log.WithField("account_address", accID).Debug("Added Account to BlackList.")
		s.blackList[accID] = struct{}{}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.requestHandler)

	var server *http.Server
	go running.UntilSuccess(ctx, s.log, "listening_server", func(ctx context.Context) (bool, error) {
		server = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", s.config.Listener.Host, s.config.Listener.Port),
			Handler:      mux,
			WriteTimeout: s.config.Listener.Timeout,
		}

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return false, errors.Wrap(err, "Failed to ListenAndServe (Server stopped with error)")
		}

		return false, nil
	}, time.Second, time.Hour)

	<-ctx.Done()

	shutdownCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	server.Shutdown(shutdownCtx)
	s.log.Info("Server stopped cleanly.")
}
