package idmind

import (
	"net/http"
	"strings"

	"context"

	"io/ioutil"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/go/xdr"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/kyc"
	"gitlab.com/swarmfund/psim/psim/templates"
)

// ProcessNotSubmitted approves Users from USA or with non-Latin document,
// submits all others KYCs to IDMind.
func (s *Service) processNotSubmitted(ctx context.Context, request horizon.Request) error {
	logger := s.log.WithField("request", request)
	kycRequest := request.Details.KYC

	accountID := kycRequest.AccountToUpdateKYC
	fields := logan.F{
		"account_id": accountID,
	}

	if kycRequest.AccountTypeToSet.Int != int(xdr.AccountTypeGeneral) {
		// Mark as reviewed without sending to IDMind (it's probably Syndicate - we don't handle Syndicates via IDMind)
		err := s.approveBothTasks(ctx, request.ID, request.Hash, false)
		if err != nil {
			return errors.Wrap(err, "Failed to approve both Tasks (without sending to IDMind - nonGeneral Account requested)")
		}

		logger.Info("Successfully approved without sending to IDMind (nonGeneral requested).")
		return nil
	}

	blob, emailAddr, err := s.retrieveBlobAndEmail(request, accountID)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve Blob or email", fields)
	}
	kycData, err := kyc.ParseKYCData(blob.Attributes.Value)
	if err != nil {
		// Blob data is unparsable - rejecting.
		_, err = s.reject(ctx, request.ID, request.Hash, nil, "Something went wrong, please try again", 0, map[string]string{
			"additional_info": "Tried to parse KYC data from Blob.Attributes.Values, but failed.",
		})
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest because of unparsable KYCData from Blob")
		}

		logger.WithField("blob", blob).Info("Successfully rejected KYCRequest (because of unparsable KYCData from Blob).")
		return nil
	}

	isUSA := kycData.IsUSA()
	if kycRequest.AllTasks&TaskNonLatinDoc != 0 || isUSA {
		// Mark as reviewed without sending to IDMind (non-Latin document or from USA - IDMind doesn't handle such guys)
		err := s.approveWithoutSubmit(ctx, request, isUSA, kycData.FirstName, emailAddr)
		if err != nil {
			return errors.Wrap(err, "Failed to approve without sending to IDMind - nonLatin docs or USA")
		}

		return nil
	}

	err = s.processNewKYCApplication(ctx, *kycData, emailAddr, accountID, request)
	if err != nil {
		return errors.Wrap(err, "Failed to process new KYC Application", fields)
	}

	return nil
}

// RetrieveKYCData retrieves BlobID from KYCRequest,
// obtains Blob by BlobID,
// check Blob's type
// and parses KYCData from the Blob.Attributes.Value
func (s *Service) retrieveBlobAndEmail(request horizon.Request, accountID string) (blob *horizon.Blob, email string, err error) {
	user, err := s.usersConnector.User(accountID)
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to get User by AccountID from Horizon")
	}
	email = user.Attributes.Email

	kycRequest := request.Details.KYC

	blobIDInterface, ok := kycRequest.KYCData["blob_id"]
	if !ok {
		return nil, email, errors.New("Cannot found 'blob_id' key in the KYCData map in the KYCRequest.")
	}

	blobID, ok := blobIDInterface.(string)
	if !ok {
		// Normally should never happen
		return nil, email, errors.New("BlobID from KYCData map of the KYCRequest is not a string.")
	}

	blob, err = s.blobsConnector.Blob(blobID)
	if err != nil {
		return nil, email, errors.Wrap(err, "Failed to get Blob from Horizon")
	}
	fields := logan.F{"blob": blob}

	if blob.Type != KYCFormBlobType {
		return nil, email, errors.From(errors.Errorf("The Blob provided in KYC Request is of type (%s), but expected (%s).",
			blob.Type, KYCFormBlobType), fields)
	}

	return blob, email, nil
}

func (s *Service) approveWithoutSubmit(ctx context.Context, request horizon.Request, isUSA bool, firstName, emailAddr string) error {
	err := s.approveBothTasks(ctx, request.ID, request.Hash, isUSA)
	if err != nil {
		return errors.Wrap(err, "Failed to approve both Tasks")
	}

	var logDetail string
	if isUSA {
		logDetail = "USA User"
	} else {
		logDetail = "nonLatin docs"
	}
	s.log.WithFields(logan.F{
		"is_usa":  isUSA,
		"request": request,
	}).Infof("Successfully approved without sending to IDMind - %s.", logDetail)

	if isUSA {
		// Notify User, he needs to pass AccreditedInvestor check
		templateData := struct {
			Link      string
			FirstName string
		}{
			Link:      s.config.TemplateLinkURL,
			FirstName: firstName,
		}

		msg, err := templates.BuildTemplateEmailMessage(s.config.USAUsersEmailConfig.Message, templateData)
		if err != nil {
			return errors.Wrap(err, "Failed to build Email message from Template", logan.F{
				"template_name": s.config.USAUsersEmailConfig.Message,
			})
		}
		s.usaUsersEmail.AddTask(ctx, emailAddr, s.config.USAUsersEmailConfig.Subject, msg)
	}

	return nil
}

func (s *Service) processNewKYCApplication(ctx context.Context, kycData kyc.Data, email, accID string, request horizon.Request) error {
	idDoc := kycData.Documents.IDDocument

	docType := getDocType(idDoc.Type)
	fields := logan.F{"doc_type": docType}

	faceFile, backFile, err := s.fetchIDDocument(idDoc)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch Documents")
	}
	fields = fields.Merge(logan.F{
		"face_doc_file_len": len(faceFile),
		"back_doc_file_len": len(backFile),
	})

	createAccountReq, validationErr := buildCreateAccountRequest(kycData, email, docType, faceFile, backFile)
	if validationErr != nil {
		err := s.rejectInvalidKYCData(ctx, request.ID, request.Hash, kycData.IsUSA(), validationErr)
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

	err = s.processNewApplicationResponse(ctx, *applicationResponse, kycData, request, email)
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

func (s *Service) processNewApplicationResponse(ctx context.Context, appResponse ApplicationResponse, kycData kyc.Data, request horizon.Request, email string) error {
	logger := s.log.WithFields(logan.F{
		"request": request,
		"email":   email,
	})

	rejectReason, details := s.getAppRespRejectReason(appResponse)
	if rejectReason != "" {
		// Need to reject
		blobID, err := s.rejectSubmitKYC(ctx, request.ID, request.Hash, appResponse, rejectReason, details, kycData.IsUSA())
		if err != nil {
			return errors.Wrap(err, "Failed to reject KYCRequest due to reason from immediate ApplicationResponse")
		}

		logger.WithFields(logan.F{
			"reject_blob_id":     blobID,
			"reject_ext_details": details,
		}).Infof("Rejected KYCRequest during Submit Task successfully (%s).", rejectReason)
		return nil
	}

	// TODO Make sure we need TxID, not MTxID
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
