package idmind

import (
	"net/http"
	"strings"

	"context"

	"io/ioutil"

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
	idDoc := kycData.Documents.IDDocument

	var docType DocType
	switch idDoc.Type {
	case kyc.PassportDocType:
		docType = PassportDocType
	case kyc.DrivingLicenseDocType:
		docType = DrivingLicenseDocType
	case kyc.IdentityCardDocType:
		docType = IdentityCardDocType
	case kyc.ResidencePermitDocType:
		docType = ResidencePermitDocType
	}

	faceFile, backFile, err := s.fetchIDDocument(idDoc)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch Documents")
	}
	fields := logan.F{
		"doc_type":          docType,
		"face_doc_file_len": len(faceFile),
		"back_doc_file_len": len(backFile),
	}

	createAccountReq, reqValidationErr := buildCreateAccountRequest(kycData, email, docType, faceFile, backFile)
	if reqValidationErr != nil {
		err := s.rejectInvalidKYCData(ctx, request.ID, request.Hash, kycData.IsUSA(), reqValidationErr)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of invalid KYCData", fields)
		}

		// This log is of level Warn intentionally, as it's not normal situation, front-end must always provide valid KYC data
		s.log.WithField("request", request).Warn("Rejected KYCRequest during Submit Task successfully (invalid KYC data).")
		return nil
	}

	applicationResponse, err := s.identityMind.Submit(*createAccountReq)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYC data to IdentityMind", fields)
	}
	fields["application_response"] = applicationResponse

	err = s.processNewApplicationResponse(ctx, *applicationResponse, kycData, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process response of new KYC Application", fields)
	}

	return nil
}

func (s *Service) fetchIDDocument(doc kyc.IDDocument) (faceFile, backFile []byte, err error) {
	faceDoc, err := s.documentsConnector.Document(doc.FaceDocID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to get Face Document by ID from Horizon")
	}

	faceFileResp, err := http.Get(fixDocURL(faceDoc.URL))
	faceFile, err = ioutil.ReadAll(faceFileResp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to read faceFile response into bytes")
	}

	if doc.BackDocID != "" {
		backDoc, err := s.documentsConnector.Document(doc.BackDocID)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to get Back Document by ID from Horizon")
		}

		backFileResp, err := http.Get(fixDocURL(backDoc.URL))
		backFile, err = ioutil.ReadAll(backFileResp.Body)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to read backFile response into bytes")
		}
	}

	return faceFile, backFile, nil
}

func fixDocURL(url string) string {
	return strings.Replace(url, `\u0026`, `&`, -1)
}

func (s *Service) processNewApplicationResponse(ctx context.Context, appResponse ApplicationResponse, kycData kyc.Data, request horizon.Request) error {
	rejectReason, details := s.getAppRespRejectReason(appResponse)
	if rejectReason != "" {
		// Need to reject
		blobID, err := s.rejectSubmitKYC(ctx, request.ID, request.Hash, appResponse, rejectReason, details, kycData.IsUSA())
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest due to reason from immediate ApplicationResponse")
		}

		s.log.WithFields(logan.F{
			"request":        request,
			"reject_blob_id": blobID,
		}).Infof("Rejected KYCRequest during Submit Task successfully (%s).", rejectReason)
		return nil
	}

	// TODO Make sure we need TxID, not MTxID
	err := s.approveSubmitKYC(ctx, request.ID, request.Hash, appResponse.TxID, kycData.IsUSA())
	if err != nil {
		return errors.Wrap(err, "Failed to approve submit part of KYCRequest")
	}

	s.log.WithField("request", request).Info("Approved KYCRequest during Submit Task successfully.")
	return nil
}

// GetAppRespRejectReason returns "", nil if no reject reasons in immediate Application response.
func (s *Service) getAppRespRejectReason(appResponse ApplicationResponse) (rejectReason string, details map[string]string) {
	rejectReason, details = s.getCheckRespRejectReason(appResponse.CheckApplicationResponse)
	if rejectReason != "" {
		return rejectReason, details
	}

	if appResponse.FraudResult == DenyFraudResult {
		return s.config.RejectReasons.FraudPolicyResultDenied, nil
	}

	return "", nil
}
