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
	kyc := request.Details.KYC

	blobIDInterface, ok := kyc.KYCData["blob_id"]
	if !ok {
		return errors.New("Cannot found 'blob_id' key in the KYCData map in the KYCRequest.")
	}
	blobID, ok := blobIDInterface.(string)
	if !ok {
		// Normally should never happen
		return errors.New("BlobID from KYCData map of the KYCRequest is not a string.")
	}

	err := s.processKYCBlob(ctx, request, blobID, kyc.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "Failed to process KYC Blob", logan.F{
			"blob_id":    blobID,
			"account_id": kyc.AccountToUpdateKYC,
		})
	}

	return nil
}

// TODO Refactor me - too long method
func (s *Service) processKYCBlob(ctx context.Context, request horizon.Request, blobID, accountID string) error {
	blob, err := s.blobProvider.Blob(blobID)
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

	user, err := s.userProvider.User(accountID)
	if err != nil {
		return errors.Wrap(err, "Failed to get User by AccountID from Horizon", fields)
	}
	email := user.Attributes.Email

	createAccountReq, err := buildCreateAccountRequest(*kycData, email)
	if err != nil {
		err := s.rejectInvalidKYCData(ctx, request.ID, request.Hash, kycData.IsUSA(), err)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of invalid KYCData", fields)
		}

		// This log is of level Warn intentionally, as it's not normal situation, front-end must always provide valid KYC data
		s.log.WithField("request", request).Warn("Rejected KYCRequest during Submit Task successfully (invalid KYC data).")
		return nil
	}

	applicationResponse, err := s.identityMind.Submit(*createAccountReq)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYC data to IdentityMind")
	}
	fields["app_response"] = applicationResponse

	if applicationResponse.KYCState == RejectedKYCState {
		err := s.rejectSubmitKYC(ctx, request.ID, request.Hash, applicationResponse, s.config.RejectReasons.KYCStateRejected, kycData.IsUSA())
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of KYCState rejected in immediate ApplicationResponse", fields)
		}

		s.log.WithField("request", request).Info("Rejected KYCRequest during Submit Task successfully (rejected state).")
		return nil
	}
	if applicationResponse.PolicyResult == DenyFraudResult {
		err := s.rejectSubmitKYC(ctx, request.ID, request.Hash, applicationResponse, s.config.RejectReasons.FraudPolicyResultDenied, kycData.IsUSA())
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of PolicyResult(fraud) denied in immediate ApplicationResponse", fields)
		}

		s.log.WithField("request", request).Info("Rejected KYCRequest during Submit Task successfully (denied FraudPolicyResult).")
		return nil
	}

	// TODO Make sure we need TxID, not MTxID
	err = s.fetchAndSubmitDocs(kycData.Documents, applicationResponse.TxID)
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

func (s *Service) fetchAndSubmitDocs(docs kyc.Documents, txID string) error {
	doc, err := s.documentProvider.Document(docs.KYCIdDocument)
	if err != nil {
		return errors.Wrap(err, "Failed to get KYCIdDocument by ID from Horizon")
	}

	resp, err := http.Get(fixDocURL(doc.URL))
	// TODO parse response Content-Type to determine document file extension (do when it's ready in API)

	err = s.identityMind.UploadDocument(txID, "ID Document", "id_document", resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYCIdDocument to IdentityMind")
	}

	doc, err = s.documentProvider.Document(docs.KYCProofOfAddress)
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
