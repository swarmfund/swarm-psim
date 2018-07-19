package bearer

import (
	"context"

	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/ape/problems"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/horizon-connector"

	"strconv"

	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/psim/psim/listener"
)

type SalesQ interface {
	CoreSales() ([]horizon.CoreSale, error)
}

// CheckSalesStateHelperInterface is an interface allows to perform check sale state operation
// and get all required data for it
type CheckSalesStateHelperInterface interface {
	SalesQ
	CloseSale(id uint64) (bool, error)
}

// Service is a main structure for bearer runner,
// implements `app.Service` interface.
type Service struct {
	config       Config
	helper       CheckSalesStateHelperInterface
	logger       *logan.Entry
	salesToClose chan horizon.CoreSale
}

// New is a constructor for bearer Service.
func New(config Config, log *logan.Entry, helper CheckSalesStateHelperInterface) *Service {
	return &Service{
		config:       config,
		helper:       helper,
		logger:       log.WithField("service", conf.ServiceBearer),
		salesToClose: make(chan horizon.CoreSale, 10),
	}
}

func (s *Service) handleCloseSaleRequest(w http.ResponseWriter, r *http.Request) {
	saleID := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(saleID, 10, 64)
	if len(saleID) == 0 || err != nil {
		ape.RenderErr(w, problems.BadRequest("invalid id"))
		return
	}
	s.salesToClose <- horizon.CoreSale{ID: id}
	w.WriteHeader(http.StatusAccepted)
}

func (s *Service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.logger),
	)

	r.Put("/sale/{id}", s.handleCloseSaleRequest)

	return r
}

// Run logs out info about service start and invokes run over incremental timer
func (s *Service) Run(ctx context.Context) {
	s.logger.Info("starting...")

	go running.WithBackOff(ctx,
		s.logger,
		"crawler",
		s.crawler,
		s.config.NormalTime,
		s.config.AbnormalPeriod,
		s.config.MaxAbnormalPeriod)

	go running.WithBackOff(ctx,
		s.logger,
		"saleCloser",
		s.saleCloser,
		s.config.NormalTime,
		s.config.AbnormalPeriod,
		s.config.MaxAbnormalPeriod)

	serverConf := listener.Config{
		Host:           s.config.Host,
		CheckSignature: false,
		Port:           s.config.Port,
	}
	listener.RunServer(ctx, s.logger, s.router(), serverConf)
}

func (s *Service) tryCloseSale(sale horizon.CoreSale) error {
	fields := logan.F{
		"sale_id": sale.ID,
	}
	closed, err := s.helper.CloseSale(sale.ID)
	if err != nil {
		return err
	}
	if closed {
		s.logger.WithFields(fields).Debug("sale closed")
	} else {
		s.logger.WithFields(fields).Debug("sale not ready yet")
	}
	return nil
}

func (s *Service) saleCloser(ctx context.Context) error {
	for {
		select {
		case sale := <-s.salesToClose:
			fields := logan.F{
				"sale_id": sale.ID,
			}
			err := s.tryCloseSale(sale)
			if err != nil {
				s.logger.WithFields(fields).WithError(err).Error("failed to close sale")
				return errors.Wrap(err, "failed to close sale", fields)
			}
			s.logger.WithFields(fields).Debug("Handled close sale request")
		case <-ctx.Done():
			s.logger.Info("saleCloser shut down")
		}
	}

	return nil
}

func (s *Service) crawler(ctx context.Context) error {
	sales, err := s.helper.CoreSales()
	if err != nil {
		s.logger.WithError(err).Error("failed to get sales")
		return errors.Wrap(err, "failed to get sales")
	}

	for _, sale := range sales {
		fields := logan.F{
			"sale_id": sale.ID,
		}
		err := s.tryCloseSale(sale)
		if err != nil {
			s.logger.WithFields(fields).WithError(err).Error("failed to close sale")
			return errors.Wrap(err, "failed to close sale", fields)
		}
	}
	return nil
}
