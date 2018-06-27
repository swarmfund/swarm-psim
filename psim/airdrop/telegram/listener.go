package telegram

import (
	"context"
	"encoding/json"
	"net/http"

	"fmt"

	"strings"

	"gitlab.com/swarmfund/psim/psim/airdrop"
	"gitlab.com/swarmfund/psim/psim/listener"
)

func (s *Service) requestHandler(w http.ResponseWriter, r *http.Request) {
	bb, errResponseWritten := listener.ValidateHTTPRequest(w, r, s.log, s.doorman)
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

	if err := request.Validate(); err != "" {
		logger.WithField("validation_err", err).Warn("Received invalid request.")
		listener.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if _, ok := s.blackList[request.AccountID]; ok {
		logger.Warn("Received issuance request to a black-listed AccountID.")
		listener.WriteError(w, http.StatusForbidden, "This Account is forbidden to receive this Airdrop.")
		return
	}

	if ok := s.checkHandle(ctx, w, request); !ok {
		return
	}

	balanceID, err := s.balanceIDProvider.GetBalanceID(request.AccountID, s.config.Issuance.Asset)
	if err != nil {
		logger.WithError(err).Error("Failed to get BalanceID.")
		listener.WriteError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}
	if balanceID == nil {
		// It can also happen if no Account was found.
		logger.Warn("No Balance was found.")
		listener.WriteError(w, http.StatusNotFound, "No Balance was found.")
		return
	}
	logger = logger.WithField("balance_id", balanceID)

	//issuanceOpt, issuanceHappened, err := s.issueSWM(ctx, request.AccountID)
	opDetails := fmt.Sprintf(`{"cause": "%s"}`, airdrop.TelegramIssuanceCause)
	issuanceOpt, issuanceHappened, err := s.issuanceSubmitter.Submit(ctx, request.AccountID, *balanceID, s.config.Issuance.Amount, opDetails)
	if err != nil {
		logger.WithError(err).Error("Failed to fulfill SWM issuance request.")
		listener.WriteError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}

	if !issuanceHappened {
		logger.WithField("issuance", issuanceOpt).Info("Reference duplication.")
		listener.WriteResponse(w, http.StatusNoContent, nil)
		return
	}

	logger.WithField("issuance", issuanceOpt).Info("New issuance happened.")
	response := fmt.Sprintf(`{"issuance_reference":"%s"}`, issuanceOpt.Reference)
	bb := []byte(response)
	listener.WriteResponse(w, http.StatusCreated, bb)
	return
}

func (s *Service) checkHandle(ctx context.Context, w http.ResponseWriter, request UserRequest) bool {
	logger := s.log.WithField("request", request)

	handle := request.TelegramHandle
	if strings.HasPrefix(handle, "@") {
		handle = handle[:len(handle)-1]
	}

	handleDoesntExist, err := s.connector.CheckUsername(ctx, handle)
	if err != nil {
		logger.WithError(err).Error("Failed to check existence of Telegram handle.")
		listener.WriteError(w, http.StatusInternalServerError, "Internal error occurred.")
		return false
	}
	if handleDoesntExist {
		logger.Info("Telegram handle does not exist.")
		listener.WriteError(w, http.StatusNotFound, fmt.Sprintf("Telegram handle '%s' does not exist.", request.TelegramHandle))
		return false
	}

	return true
}
