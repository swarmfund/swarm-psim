package verification

import (
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/swarmfund/psim/ape"
	"gitlab.com/swarmfund/psim/ape/problems"
)

// TODO Consider moving to verifier package if withdraw moves to common verifier eventually.
func ReadAPIRequest(log *logan.Entry, w http.ResponseWriter, r *http.Request, request interface{}) (success bool) {
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.WithError(err).Warn("Failed to parse request.")
		ape.RenderErr(w, r, problems.BadRequest("Cannot parse JSON request."))
		return false
	}

	return true
}

// TODO Start rendering only signature in the response, not the whole signed Envelope.
func RenderResponseEnvelope(log *logan.Entry, w http.ResponseWriter, r *http.Request, envelope string) (success bool) {
	response := VerifyResponse{
		Envelope: envelope,
	}

	respBytes, err := json.Marshal(response)
	if err != nil {
		log.WithField("response_trying_to_render", response).WithError(err).Error("Failed to marshal VerifyResponse.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return false
	}

	_, err = w.Write(respBytes)
	if err != nil {
		log.WithField("envelope_response", string(respBytes)).WithError(err).
			Error("Failed to write VerifyResponse bytes into the ResponseWriter.")
		ape.RenderErr(w, r, problems.ServerError(err))
		return false
	}

	return true
}
