package idmind

import (
	"strings"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
)

func (s *Service) submitKYCToIDMind(request horizon.Request) error {
	kyc := request.Details.KYC

	blobIDInterface, ok := kyc.KYCData["blob_id"]
	if !ok {
		return errors.New("Cannot found 'blob_id' key in the KYCData map in the KYCRequest.")
	}
	blobID, ok := blobIDInterface.(string)
	if !ok {
		return errors.New("BlobID from KYCData map of the KYCRequest is not a string.")
	}

	err := s.submitKYCBlob(blobID, kyc.AccountToUpdateKYC)
	if err != nil {
		return errors.Wrap(err, "Failed to process KYC Blob", logan.F{
			"blob_id":    blobID,
			"account_id": kyc.AccountToUpdateKYC,
		})
	}

	return nil
}

// TODO
func (s *Service) submitKYCBlob(blobID string, accountID string) error {
	blob, err := s.blobProvider.Blob(blobID)
	if err != nil {
		return errors.Wrap(err, "Failed to get Blob from Horizon")
	}
	fields := logan.F{"blob": blob}

	if blob.Type != KYCFormBlobType {
		return errors.From(errors.Errorf("The Blob provided in KYC Request is of type (%s), but expected (%s).",
			blob.Type, KYCFormBlobType), fields)
	}

	kycData, err := parseKYCData(blob.Attributes.Value)
	if err != nil {
		return errors.Wrap(err, "Failed to parse KYC data from Attributes.Value string in from Blob", fields)
	}
	fields["kyc_data"] = kycData

	user, err := s.userProvider.User(accountID)
	if err != nil {
		return errors.Wrap(err, "Failed to get User by AccountID from Horizon", fields)
	}
	email := user.Attributes.Email

	applicationResponse, err := s.identityMind.Submit(*kycData, email)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYC data to IdentityMind")
	}

	// TODO
	if applicationResponse.KYCState == RejectedKYCState {
		// TODO Reject KYC request with specific RejectReason
	}
	if applicationResponse.PolicyResult == DenyFraudResult {
		// TODO Reject KYC request with specific RejectReason
	}

	// TODO Make sure we need TxID, not MTxID
	err = s.fetchAndSubmitDocs(kycData.Documents, applicationResponse.TxID)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch and submit KYC documents")
	}

	// TODO Updated KYC ReviewableRequest with TxID, response-result?, ... got from IM (submit Op)

	return errors.New("Not implemented.")
}

func (s *Service) fetchAndSubmitDocs(docs KYCDocuments, txID string) error {
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
