package listener

import (
	"github.com/pkg/errors"
	"gitlab.com/swarmfund/psim/salesforce"
)

type SalesforceConnector struct {
	*salesforce.Connector
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

func NewSalesforceConnector(connector *salesforce.Connector) *SalesforceConnector {
	return &SalesforceConnector{
		connector,
	}
}

// SendEvent tries to send an event, if auth token expired - get new, if error - retry
func (sc *SalesforceConnector) SendEvent(event BroadcastedEvent) error {
	_, err := sc.Connector.SendEvent(eventNameToSphere[event.Name], eventNameToActionName[event.Name], event.Time.Format(salesforceTimeLayout), event.ActorName, event.ActorEmail, event.InvestmentAmount, event.InvestmentCountry)
	if err != nil {
		return errors.Wrap(err, "failed to post event")
	}
	// TODO handle resp
	return nil
}
