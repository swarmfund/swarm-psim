package idmind

import (
	"net/http"
	"strings"

	"context"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/kyc"
)

func (s *Service) submitKYCData(ctx context.Context, request horizon.Request) error {
	kycRequest := request.Details.KYC

	blobIDInterface, ok := kycRequest.KYCData["blob_id"]
	if !ok {
		return errors.New("Cannot found 'blob_id' key in the KYCData map in the KYCRequest.")
	}

	blobID, ok := blobIDInterface.(string)
	if !ok {
		// Normally should never happen
		return errors.New("BlobID from KYCData map of the KYCRequest is not a string.")
	}

	err := s.processKYCBlob(ctx, request, blobID, kycRequest.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "Failed to process KYC Blob", logan.F{
			"blob_id":    blobID,
			"account_id": kycRequest.AccountToUpdateKYC,
		})
	}

	return nil
}

func (s *Service) processKYCBlob(ctx context.Context, request horizon.Request, blobID, accountID string) error {
	blob, err := s.blobsConnector.Blob(blobID)
	if err != nil {
		return errors.Wrap(err, "Failed to get Blob from Horizon")
	}
	fields := logan.F{"blob": blob}

	if blob.Type != KYCFormBlobType {
		return errors.From(errors.Errorf("The Blob provided in KYC Request is of type (%s), but expected (%s).",
			blob.Type, KYCFormBlobType), fields)
	}

	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		return errors.Wrap(err, "Failed to parse KYC data from Attributes.Value string in from Blob", fields)
	}
	fields["kyc_data"] = kycData

	user, err := s.usersConnector.User(accountID)
	if err != nil {
		return errors.Wrap(err, "Failed to get User by AccountID from Horizon", fields)
	}
	email := user.Attributes.Email

	err = s.processNewKYCApplication(ctx, *kycData, email, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process new KYC Application", fields)
	}

	return nil
}

func (s *Service) processNewKYCApplication(ctx context.Context, kycData kyc.Data, email string, request horizon.Request) error {
	createAccountReq, err := buildCreateAccountRequest(kycData, email)
	if err != nil {
		err := s.rejectInvalidKYCData(ctx, request.ID, request.Hash, kycData.IsUSA(), err)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of invalid KYCData")
		}

		// This log is of level Warn intentionally, as it's not normal situation, front-end must always provide valid KYC data
		s.log.WithField("request", request).Warn("Rejected KYCRequest during Submit Task successfully (invalid KYC data).")
		return nil
	}

	applicationResponse, err := s.identityMind.Submit(*createAccountReq)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYC data to IdentityMind")
	}

	err = s.processNewApplicationResponse(ctx, *applicationResponse, kycData, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process response of new KYC Application", logan.F{
			"app_response": applicationResponse,
		})
	}

	return nil
}

func (s *Service) processNewApplicationResponse(ctx context.Context, applicationResponse ApplicationResponse, kycData kyc.Data, request horizon.Request) error {
	if applicationResponse.KYCState == RejectedKYCState {
		blobID, err := s.rejectSubmitKYC(ctx, request.ID, request.Hash, applicationResponse, s.config.RejectReasons.KYCStateRejected, kycData.IsUSA())
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of KYCState rejected in immediate ApplicationResponse")
		}

		s.log.WithFields(logan.F{
			"request":        request,
			"reject_blob_id": blobID,
		}).Info("Rejected KYCRequest during Submit Task successfully (rejected state).")
		return nil
	}
	if applicationResponse.PolicyResult == DenyFraudResult {
		blobID, err := s.rejectSubmitKYC(ctx, request.ID, request.Hash, applicationResponse, s.config.RejectReasons.FraudPolicyResultDenied, kycData.IsUSA())
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of PolicyResult(fraud) denied in immediate ApplicationResponse")
		}

		s.log.WithFields(logan.F{
			"request":        request,
			"reject_blob_id": blobID,
		}).Info("Rejected KYCRequest during Submit Task successfully (denied FraudPolicyResult).")
		return nil
	}

	// TODO Make sure we need TxID, not MTxID
	err := s.fetchAndSubmitDocs(kycData.Documents, applicationResponse.TxID)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch and submit KYC documents")
	}

	// TODO Make sure we need TxID, not MTxID
	err = s.approveSubmitKYC(ctx, request.ID, request.Hash, applicationResponse.TxID, kycData.IsUSA())
	if err != nil {
		return errors.Wrap(err, "Failed to approve submit part of KYCRequest")
	}

	s.log.WithField("request", request).Info("Approved KYCRequest during Submit Task successfully.")
	return nil
}

// FIXME
func (s *Service) fetchAndSubmitDocs(docs kyc.Documents, txID string) error {
	doc, err := s.documentsConnector.Document(docs.KYCIdDocument)
	if err != nil {
		return errors.Wrap(err, "Failed to get KYCIdDocument by ID from Horizon")
	}

	resp, err := http.Get(fixDocURL(doc.URL))
	// TODO parse response Content-Type to determine document file extension (do when it's ready in API)

	err = s.identityMind.UploadDocument(txID, "ID Document", "id_document", resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYCIdDocument to IdentityMind")
	}

	doc, err = s.documentsConnector.Document(docs.KYCProofOfAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to get KYCProofOfAddress by ID from Horizon")
	}

	resp, err = http.Get(fixDocURL(doc.URL))
	// TODO parse response Content-Type to determine document file extension (do when it's ready in API)

	err = s.identityMind.UploadDocument(txID, "Proof of Address", "proof_of_address", resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYCProofOfAddress document to IdentityMind")
	}

	return nil
}

func fixDocURL(url string) string {
	return strings.Replace(url, `\u0026`, `&`, -1)
}
