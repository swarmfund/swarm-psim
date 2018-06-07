package listener

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

// SalesforceTarget represents a target with salesforce api
type SalesforceTarget struct {
	SalesforceClient
	Username     string
	Password     string
	ClientSecret string
	AccessToken  string
	APIUrl       *url.URL
}

const emptyAccessToken = ""

func (st *SalesforceTarget) auth() error {
	code, resp, err := st.PostAuthRequest(st.Username, st.Password, st.ClientSecret)
	if err != nil {
		return errors.Wrap(err, "post auth request using salesforce client failed")
	}
	if code != 200 {
		return errors.New("got unsuccessful http code")
	}
	authResponse := SalesforceAuthResponse{}
	err = json.Unmarshal(resp, &authResponse)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response")
	}
	st.AccessToken = authResponse.AccessToken
	return nil
}

// NewSalesforceTarget constructs a target
func NewSalesforceTarget(username string, password string, clientSecret string, APIUrl string) (*SalesforceTarget, error) {
	targetURL, err := url.Parse(APIUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse salesforce api targetURL")
	}

	httpClient := &http.Client{}

	salesforceClient := NewSalesforceClient(httpClient, targetURL)

	target := &SalesforceTarget{*salesforceClient, username, password, clientSecret, emptyAccessToken, targetURL}

	err = target.auth()
	if err != nil {
		return nil, errors.Wrap(err, "failed to auth while creating salesforce target")
	}

	return target, nil
}

const salesforceTimeLayout = "2006-01-02T15:04:05.999-0700"

var eventNameToSphere = map[BroadcastedEventName]string{
	BroadcastedEventNameKycCreated:            "Total KYC",
	BroadcastedEventNameKycUpdated:            "Total KYC",
	BroadcastedEventNameKycRejected:           "Total KYC",
	BroadcastedEventNameKycApproved:           "Total KYC",
	BroadcastedEventNameUserReferred:          "Community",
	BroadcastedEventNameFundsWithdrawn:        "Investment",
	BroadcastedEventNamePaymentV2Received:     "Trading",
	BroadcastedEventNamePaymentV2Sent:         "Trading",
	BroadcastedEventNamePaymentReceived:       "Trading",
	BroadcastedEventNamePaymentSent:           "Trading",
	BroadcastedEventNameFundsDeposited:        "Investment",
	BroadcastedEventNameFundsInvested:         "Investment",
	BroadcastedEventNameReferredUserPassedKyc: "Community",
}

var eventNameToActionName = map[BroadcastedEventName]string{
	BroadcastedEventNameKycCreated:            "Submit KYC",
	BroadcastedEventNameKycUpdated:            "Resubmit KYC",
	BroadcastedEventNameKycRejected:           "KYC Rejected",
	BroadcastedEventNameKycApproved:           "Complete KYC",
	BroadcastedEventNameUserReferred:          "Refer a user",
	BroadcastedEventNameFundsWithdrawn:        "Withdraw Funds",
	BroadcastedEventNamePaymentV2Received:     "Receive Funds",
	BroadcastedEventNamePaymentV2Sent:         "Send Funds",
	BroadcastedEventNamePaymentReceived:       "Receive Funds",
	BroadcastedEventNamePaymentSent:           "Send Funds",
	BroadcastedEventNameFundsDeposited:        "Deposit Funds",
	BroadcastedEventNameFundsInvested:         "Invest Funds",
	BroadcastedEventNameReferredUserPassedKyc: "Referred user completed KYC",
}

// SendEvent tries to send an event, if auth token expired - get new, if error - retry
func (st *SalesforceTarget) SendEvent(event BroadcastedEvent) error {
	statusCode, resp, err := st.PostEvent(st.AccessToken, eventNameToSphere[event.Name], eventNameToActionName[event.Name], event.Time.Format(salesforceTimeLayout), event.ActorName, event.ActorEmail, event.InvestmentAmount, event.InvestmentCountry)

	if err != nil {
		return errors.Wrap(err, "failed to post event")
	}

	if statusCode != 201 {
		st.auth()
		return errors.New("got " + fmt.Sprint(statusCode) + ", reauth performed")
	}

	respStruct := SalesforceEventResponse{}
	err = json.Unmarshal(resp, &respStruct)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal response while sending event")
	}

	return nil
}
