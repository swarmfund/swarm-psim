package investready

import (
	"net/http"
	"context"
	"encoding/json"
)

func (l *RedirectsListener) userHashHandler(w http.ResponseWriter, r *http.Request) {
	bb, errResponseWritten := l.validateHTTPRequest(w, r, http.MethodPut)
	if errResponseWritten {
		return
	}

	var request userHashRequest
	err := json.Unmarshal(bb, &request)
	if err != nil {
		l.log.WithField("raw_request", string(bb)).WithError(err).Warn("Failed to unmarshal UserHash request bytes into struct.")
		writeError(w, http.StatusBadRequest, "Cannot parse JSON request.")
		return
	}

	l.processUserHashRequest(r.Context(), w, request)
}

func (l *RedirectsListener) processUserHashRequest(ctx context.Context, w http.ResponseWriter, request userHashRequest) {
	logger := l.log.WithField("request", request)

	if validationErr := request.Validate(); validationErr != "" {
		logger.WithField("validation_err", validationErr).Warn("Received invalid request.")
		writeError(w, http.StatusBadRequest, validationErr)
		return
	}

	kycRequest, forbiddenErr, err := l.getKYCRequest(ctx, request.AccountID)
	if err != nil {
		logger.WithError(err).Error("Failed to get KYCRequest by AccountID.")
		writeError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}
	if forbiddenErr != nil {
		logger.WithField("forbidden_reason", forbiddenErr).Warn("User is forbidden to add InvestReady UserHash to the KYCRequest.")
		writeError(w, http.StatusForbidden, forbiddenErr.Error())
		return
	}

	l.saveUserHash(ctx, *kycRequest, request.AccountID, request.UserHash)
	if err != nil {
		logger.WithError(err).Error("Failed to save UserHash.")
		writeError(w, http.StatusInternalServerError, "Internal error occurred.")
		return
	}

	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.WriteHeader(http.StatusNoContent)
}
