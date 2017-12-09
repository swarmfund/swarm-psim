package stripeverify

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/stripe/stripe-go/client"
	"gitlab.com/distributed_lab/discovery-go"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/figure"
	"gitlab.com/swarmfund/psim/psim/app"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/swarmfund/psim/psim/utils"
)

func init() {
	app.Register(conf.ServiceStripeVerify, func(ctx context.Context) error {
		serviceConfig := Config{
			Host:        "localhost",
			ServiceName: conf.ServiceStripeVerify,
		}
		globalConfig := app.Config(ctx)
		err := figure.
			Out(&serviceConfig).
			From(globalConfig.Get(conf.ServiceStripeVerify)).
			With(figure.BaseHooks, utils.CommonHooks).
			Please()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to figure out %s", conf.ServiceStripeVerify))
		}

		log := ctx.Value(app.CtxLog).(*logan.Entry)

		discovery, err := globalConfig.Discovery()
		if err != nil {
			return errors.Wrap(err, "failed to get discovery client")
		}

		horizon, err := globalConfig.Horizon()
		if err != nil {
			return errors.Wrap(err, "failed to get horizon client")
		}

		listener, err := ape.Listener(serviceConfig.Host, serviceConfig.Port)
		if err != nil {
			return errors.Wrap(err, "failed to init listener")
		}

		stripe, err := globalConfig.Stripe()
		if err != nil {
			return errors.Wrap(err, "failed to init stripe client")
		}

		if serviceConfig.Signer == nil {
			panic("StribeVerify must have signer")
		}

		log.Info("starting")
		retryTicker := time.NewTicker(5 * time.Second)
		service, errs := New(serviceConfig, log, discovery, listener, horizon, stripe)
		for service.Run(); ; <-retryTicker.C {
			log.WithError(<-errs).Warn("service error")
		}
	})
}

type Service struct {
	ServiceID string

	config           Config
	log              *logan.Entry
	discovery        *discovery.Client
	listener         net.Listener
	horizon          *horizon.Connector
	stripe           *client.API
	errors           chan error
	discoveryService *discovery.Service
}

func New(
	config Config, log *logan.Entry, discovery *discovery.Client, listener net.Listener,
	horizon *horizon.Connector, stripe *client.API,
) (*Service, chan error) {
	service := &Service{
		ServiceID: utils.GenerateToken(),
		config:    config,
		log:       log,
		discovery: discovery,
		listener:  listener,
		horizon:   horizon,
		stripe:    stripe,
		errors:    make(chan error),
	}

	return service, service.errors
}

func (s *Service) Run() {
	seq := []func(){
		s.Register,
		s.API,
	}

	for _, fn := range seq {
		go fn()
	}
}

func (s *Service) Register() {
	s.discoveryService = s.discovery.Service(&discovery.ServiceRegistration{
		Name: s.config.ServiceName,
		ID:   s.ServiceID,
		Host: fmt.Sprintf("http://%s", s.listener.Addr().String()),
	})
	ticker := time.NewTicker(5 * time.Second)
	for ; true; <-ticker.C {
		err := s.discovery.RegisterServiceSync(s.discoveryService)
		if err != nil {
			s.errors <- errors.Wrap(err, "discovery error")
			continue
		}
	}
}

func (s *Service) API() {
	r := ape.DefaultRouter()

	r.Post("/", s.VerifyHandler)
	if s.config.Pprof {
		s.log.Info("enabling debugging endpoints")
		ape.InjectPprof(r)
	}

	s.log.WithField("address", s.listener.Addr().String()).Info("listening")
	s.errors <- http.Serve(s.listener, r)
}

type VerifyRequest struct {
	Envelope string `json:"envelope"`
	ChargeID string `json:"charge_id"`
}

func (s *Service) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	payload := VerifyRequest{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	transaction := s.horizon.Transaction(&horizon.TransactionBuilder{
		Envelope: payload.Envelope,
	})

	if len(transaction.Operations) != 1 {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	op := transaction.Operations[0]
	if op.Body.Type != xdr.OperationTypeManageCoinsEmissionRequest {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	body := transaction.Operations[0].Body.ManageCoinsEmissionRequestOp
	charge, err := s.stripe.Charges.Get(payload.ChargeID, nil)
	if err != nil {
		s.log.WithError(err).Warn("failed to get charge")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
	asset := charge.Meta["asset"]
	reference := charge.Meta["reference"]
	receiver := charge.Meta["receiver"]

	if asset == "" || string(body.Asset) != asset {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	if body.Action != xdr.ManageCoinsEmissionRequestActionManageCoinsEmissionRequestCreate {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	if int64(body.Amount) != int64(charge.Amount*100) {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	if receiver == "" || body.Receiver.AsString() != receiver {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	if reference == "" || string(body.Reference) != reference {
		ape.RenderErr(w, r, problems.BadRequest(""))
		return
	}

	err = transaction.Sign(s.config.Signer).Submit()
	if err != nil {
		entry := s.log.WithError(err)
		if serr, ok := err.(horizon.SubmitError); ok {
			opCodes := serr.OperationCodes()
			entry = entry.
				WithField("tx code", serr.TransactionCode()).
				WithField("op codes", serr.OperationCodes())
			if len(opCodes) == 1 {
				switch opCodes[0] {
				// safe to move on errors
				case "op_balance_not_found", "reference_duplication":
					entry.Info("tx failed")
					return
				case "op_counterparty_wrong_type":
					entry.Error("tx failed")
					return
				}
			}
		}
		entry.Error("tx failed")
		ape.RenderErr(w, r, problems.ServerError(err))
		return
	}
}
