package telegram

import (
	"context"
	"encoding/json"
	"net/http"

	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/issuance"
	"gitlab.com/swarmfund/psim/psim/listener"
)

var (
	errNoBalance = errors.New("No Balance was found.")
)

func (s *Service) requestHandler(w http.ResponseWriter, r *http.Request) {
	bb, errResponseWritten := listener.ValidateHTTPRequest(w, r, s.log, http.MethodPost, s.doorman)
	if errResponseWritten {
		return
	}

	var request UserRequest
	err := json.Unmarshal(bb, &request)
	if err != nil {
		s.log.WithField("raw_request", string(bb)).WithError(err).Warn("Failed to unmarshal UserRequest bytes into struct.")
		listener.WriteError(w, http.StatusBadRequest, "Cannot parse JSON request.")
		return
	}

	s.processUserRequest(r.Context(), w, request)
}

func (s *Service) processUserRequest(ctx context.Context, w http.ResponseWriter, request UserRequest) {
	logger := s.log.WithField("request", request)

	if validationErr := request.Validate(); validationErr != "" {
		logger.WithField("validation_err", validationErr).Warn("Received invalid request.")
		listener.WriteError(w, http.StatusBadRequest, validationErr)
		return
	}

	if _, ok := s.blackList[request.AccountID]; ok {
		logger.Warn("Received issuance request to a black-listed AccountID.")
		listener.WriteError(w, http.StatusForbidden, "This Account is forbidden to receive this Airdrop.")
		return
	}

	handleDoesntExist, err := s.connector.CheckUsername(ctx, request.TelegramHandle)
	if err != nil {
		logger.WithError(err).Error("Failed to check existence of Telegram handle.")
		listener.WriteError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}
	if handleDoesntExist {
		logger.Info("Telegram handle does not exist.")
		listener.WriteError(w, http.StatusNotFound, fmt.Sprintf("Telegram handle '%s' does not exist.", request.TelegramHandle))
		return
	}

	issuanceOpt, issuanceHappened, err := s.issueSWM(ctx, request.AccountID)
	if err != nil {
		if err == errNoBalance {
			logger.Warn("No Balance was found.")
			listener.WriteError(w, http.StatusNotFound, "No Balance was found.")
			return
		}

		logger.WithError(err).Error("Failed to fulfill SWM issuance request.")
		listener.WriteError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}

	if issuanceHappened {
		logger.WithField("issuance", issuanceOpt).Info("New issuance happened.")
		response := fmt.Sprintf(`{"issuance_reference":"%s"}`, issuanceOpt.Reference)
		bb := []byte(response)
		listener.WriteResponse(w, http.StatusCreated, bb)
		return
	} else {
		logger.WithField("issuance", issuanceOpt).Info("Reference duplication.")
		listener.WriteResponse(w, http.StatusNoContent, nil)
		return
	}
}

func (s *Service) issueSWM(ctx context.Context, accountID string) (*issuance.RequestOpt, bool, error) {
	balanceID, err := s.balanceIDProvider.GetBalanceID(accountID, s.config.Issuance.Asset)
	if err != nil {
		return nil, false, errors.Wrap(err, "Failed to get BalanceID of the Account")
	}
	if balanceID == nil {
		// It can also happen if no Account was found.
		return nil, false, errNoBalance
	}
	fields := logan.F{"balance_id": balanceID}

	opDetails := fmt.Sprintf(`{"cause": "%s"}`, airdrop.TelegramIssuanceCause)
	issuanceOpt, ok, err := s.issuanceSubmitter.Submit(ctx, accountID, *balanceID, s.config.Issuance.Amount, opDetails)
	if err != nil {
		return nil, false, errors.Wrap(err, "Failed to process Issuance", fields)
	}

	return issuanceOpt, ok, nil
}
