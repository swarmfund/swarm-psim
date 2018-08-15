package idmind

import (
	"net/http"
	"strings"

	"context"

	"io/ioutil"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/tokend/regources"
)

func (s *Service) processNewKYCApplication(ctx context.Context, kycData *kyc.Data, request regources.ReviewableRequest) error {
	idDoc := kycData.Documents.IDDocument
	docType := getDocType(idDoc.Type)
	accountID := request.Details.KYC.AccountToUpdateKYC
	fields := logan.F{
		"account_id": accountID,
		"doc_id":     idDoc,
		"doc_type":   docType,
	}

	faceFile, backFile, err := s.fetchIDDocument(idDoc)
	if err != nil {
		return errors.Wrap(err, "failed to fetch documents", fields)
	}

	fields = fields.Merge(logan.F{
		"face_doc_file_len": len(faceFile),
		"back_doc_file_len": len(backFile),
	})

	user, err := s.usersConnector.User(accountID)
	if err != nil {
		return errors.Wrap(err, "failed to get user", fields)
	}

	createAccountReq, err := buildCreateAccountRequest(
		kycData, user.Attributes.Email, user.Attributes.LastIPAddr,
		docType, faceFile, backFile,
	)
	if err != nil {
		err := s.rejectInvalidKYCData(ctx, request, err)
		if err != nil {
			return errors.Wrap(err, "Failed to reject (because of invalid KYCData) KYCRequest", fields)
		}

		// This log is of level Warn intentionally, as it's not normal situation, front-end must always provide valid KYC data
		s.log.WithField("request", request).Warn("Successfully rejected KYCRequest during Submit Task (invalid KYC data).")
		return nil
	}

	applicationResponse, err := s.identityMind.Submit(*createAccountReq)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYC data to IdentityMind", fields)
	}
	fields["application_response"] = applicationResponse

	err = s.processNewApplicationResponse(ctx, *applicationResponse, kycData, request, user.Attributes.Email)
	if err != nil {
		return errors.Wrap(err, "Failed to process response of new KYC Application", fields)
	}

	return nil
}

func getDocType(kycDocType kyc.DocType) DocType {
	switch kycDocType {
	case kyc.PassportDocType:
		return PassportDocType
	case kyc.DrivingLicenseDocType:
		return DrivingLicenseDocType
	case kyc.IdentityCardDocType:
		return IdentityCardDocType
	case kyc.ResidencePermitDocType:
		return ResidencePermitDocType
	}

	var emptyDocType DocType
	return emptyDocType
}

func (s *Service) fetchIDDocument(doc kyc.IDDocument) (faceFile, backFile []byte, err error) {
	faceDoc, err := s.documentsConnector.Document(doc.FaceFile.ID)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to get Face Document by ID from Horizon")
	}

	faceFileResp, err := http.Get(fixDocURL(faceDoc.URL))
	faceFile, err = ioutil.ReadAll(faceFileResp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to read faceFile response into bytes")
	}

	if doc.BackFile != nil {
		backDoc, err := s.documentsConnector.Document(doc.BackFile.ID)
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

// ProcessNewApplicationResponse takes email just to log it with some info logs.
func (s *Service) processNewApplicationResponse(ctx context.Context, appResponse ApplicationResponse, kycData *kyc.Data, request regources.ReviewableRequest, email string) error {
	logger := s.log.WithFields(logan.F{
		"request": request,
		"email":   email,
	})

	rejectReason, details := s.getAppRespRejectReason(appResponse)
	if rejectReason != "" {
		// Need to reject
		blobID, err := s.reject(ctx, request, appResponse, rejectReason, details)
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest due to reason from immediate ApplicationResponse")
		}

		logger.WithFields(logan.F{
			"reject_blob_id":     blobID,
			"reject_ext_details": details,
		}).Infof("Rejected KYCRequest during Submit Task successfully (%s).", rejectReason)
		return nil
	}

	err := s.approveSubmitKYC(ctx, request.ID, request.Hash, appResponse.TxID)
	if err != nil {
		return errors.Wrap(err, "Failed to approve submit part of KYCRequest")
	}
	logger = logger.WithField("tx_id", appResponse.TxID)

	logger.Info("Approved KYCRequest during Submit Task successfully.")
	return nil
}

// GetAppRespRejectReason returns "", nil if no reject reasons in immediate Application response.
func (s *Service) getAppRespRejectReason(appResponse ApplicationResponse) (rejectReason string, details map[string]string) {
	if appResponse.CheckApplicationResponse.KYCState == RejectedKYCState {
		return s.config.RejectReasons.KYCStateRejected, nil
	}

	if appResponse.FraudResult == DenyFraudResult {
		return s.config.RejectReasons.FraudPolicyResultDenied, nil
	}

	return "", nil
}
