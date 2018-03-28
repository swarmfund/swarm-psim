package idmind

import (
	"context"
	"time"

	"io"

	"strings"

	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/swarmfund/go/xdrbuild"
	"gitlab.com/swarmfund/horizon-connector/v2"
	"gitlab.com/swarmfund/psim/psim/conf"
	"gitlab.com/tokend/keypair"
)

const (
	KYCFormBlobType = "kyc_form"
)

// RequestListener is the interface, which must be implemented
// by streamer of Horizon Requests, which parametrize Service.
type RequestListener interface {
	StreamAllKYCRequests(ctx context.Context, endlessly bool) <-chan horizon.ReviewableRequestEvent
	StreamKYCRequestsUpdatedAfter(ctx context.Context, updatedAfter time.Time, endlessly bool) <-chan horizon.ReviewableRequestEvent
}

type BlobProvider interface {
	Blob(blobID string) (*horizon.Blob, error)
}

type DocumentProvider interface {
	Document(docID string) (*horizon.Document, error)
}

type UserProvider interface {
	User(accountID string) (*horizon.User, error)
}

type IdentityMind interface {
	Submit(data KYCData, email string) (*ApplicationResponse, error)
	UploadDocument(appID, txID, description string, fileName string, fileReader io.Reader) error
}

type Service struct {
	log    *logan.Entry
	config Config
	signer keypair.Full
	source keypair.Address

	requestListener  RequestListener
	blobProvider     BlobProvider
	documentProvider DocumentProvider
	userProvider     UserProvider
	identityMind     IdentityMind
	xdrbuilder       *xdrbuild.Builder

	kycRequests <-chan horizon.ReviewableRequestEvent
}

// NewService is constructor for Service.
func NewService(
	log *logan.Entry,
	config Config,
	requestListener RequestListener,
	blobProvider BlobProvider,
	userProvider UserProvider,
	documentProvider DocumentProvider,
	identityMind IdentityMind,
	builder *xdrbuild.Builder,
) *Service {

	return &Service{
		log:    log.WithField("service", conf.ServiceIdentityMind),
		config: config,

		requestListener:  requestListener,
		blobProvider:     blobProvider,
		userProvider:     userProvider,
		documentProvider: documentProvider,
		identityMind:     identityMind,
		xdrbuilder:       builder,
	}
}

// Run is a blocking method, it returns only when ctx closes.
func (s *Service) Run(ctx context.Context) {
	s.log.WithField("", s.config).Info("Starting.")

	s.kycRequests = s.requestListener.StreamAllKYCRequests(ctx, false)

	running.WithBackOff(ctx, s.log, "request_processor", s.listenAndProcessRequests, 0, 5*time.Second, 5*time.Minute)

	//appResp, err := s.identityMind.Submit(KYCData{
	//	FirstName: "John",
	//	LastName:  "Doe",
	//	Address: KYCAddress{
	//		Line1:      "Baker street",
	//		Line2:      "2B",
	//		City:       "London",
	//		Country:    "UK",
	//		State:      "CoolState",
	//		PostalCode: "123456",
	//	},
	//	ETHAddress: "",
	//	KYCDocuments:  KYCDocuments{},
	//}, "john.doe@example.com")
	//if err != nil {
	//	s.log.WithError(err).Error("Failed to submit KYC to IDMind.")
	//	return
	//}
	//s.log.WithField("app_response", appResp).Info("Received.")
}

func (s *Service) listenAndProcessRequests(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	case reqEvent := <-s.kycRequests:
		request, err := reqEvent.Unwrap()
		if err != nil {
			return errors.Wrap(err, "RequestListener sent error")
		}

		err = s.processRequest(*request)
		if err != nil {
			return errors.Wrap(err, "Failed to process KYC Request", logan.F{
				"request": request,
			})
		}

		return nil
	}
}

func (s *Service) processRequest(request horizon.Request) error {
	proveErr := proveInterestingRequest(request)
	if proveErr != nil {
		// No need to process the Request for now.
		s.log.WithError(proveErr).WithFields(logan.F{
			"request_id": request.ID,
		}).Debug("Found not interesting Request.")
		return nil
	}

	s.log.WithField("request", request).Debug("Found pending KYC Request.")

	// FIXME Take from Request
	blobID := "myAwesomeBlobID"
	// FIXME Take from Request
	accountID := "myAwesomeAccountID"

	err := s.processKYCBlob(blobID, accountID)
	if err != nil {
		return errors.Wrap(err, "Failed to process KYC Blob", logan.F{
			"blob_id":    blobID,
			"account_id": accountID,
		})
	}

	return nil
}

func (s *Service) processKYCBlob(blobID string, accountID string) error {
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

	// TODO Make user we need TxID, not MTxID
	err = s.fetchAndSubmitDocs(kycData.Documents, applicationResponse.TxID)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch and submit KYC documents")
	}

	// TODO Updated KYC ReviewableRequest with TxID, response-result?, ... got from IM (submit Op)
	// FIXME Remove this hack
	applicationResponse = applicationResponse

	return errors.New("Not implemented.")
}

func (s *Service) fetchAndSubmitDocs(docs KYCDocuments, txID string) error {
	doc, err := s.documentProvider.Document(docs.KYCIdDocument)
	if err != nil {
		return errors.Wrap(err, "Failed to get KYCIdDocument by ID from Horizon")
	}

	resp, err := http.Get(fixDocURL(doc.URL))
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" {
		// TODO ?
	}

	// FIXME appID
	// FIXME file extension
	err = s.identityMind.UploadDocument("424284", txID, "ID Document", "id_document.png", resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYCIdDocument to IdentityMind")
	}

	doc, err = s.documentProvider.Document(docs.KYCProofOfAddress)
	if err != nil {
		return errors.Wrap(err, "Failed to get KYCProofOfAddress by ID from Horizon")
	}

	resp, err = http.Get(fixDocURL(doc.URL))
	contentType = resp.Header.Get("Content-Type")
	if contentType != "" {
		// TODO ?
	}

	// FIXME appID
	// FIXME file extension
	err = s.identityMind.UploadDocument("424284", txID, "Proof of Address", "proof_of_address.png", resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to submit KYCProofOfAddress document to IdentityMind")
	}

	return nil
}

func fixDocURL(url string) string {
	return strings.Replace(url, `\u0026`, `&`, -1)
}
