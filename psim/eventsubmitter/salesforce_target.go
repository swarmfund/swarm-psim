package eventsubmitter

import (
	"github.com/pkg/errors"
	salesforce "gitlab.com/swarmfund/salesforce-connector"
)

const salesforceTimeLayout = "2006-01-02T15:04:05.999-0700"

var eventNameToSphere = map[BroadcastedEventName]string{
	BroadcastedEventNameKycCreated:            "Compilance",
	BroadcastedEventNameKycUpdated:            "Compilance",
	BroadcastedEventNameKycRejected:           "Compilance",
	BroadcastedEventNameKycApproved:           "Compilance",
	BroadcastedEventNameUserReferred:          "Community",
	BroadcastedEventNameFundsWithdrawn:        "Investment",
	BroadcastedEventNamePaymentV2Received:     "Trading",
	BroadcastedEventNamePaymentV2Sent:         "Trading",
	BroadcastedEventNamePaymentReceived:       "Trading",
	BroadcastedEventNamePaymentSent:           "Trading",
	BroadcastedEventNameFundsDeposited:        "Investment",
	BroadcastedEventNameFundsInvested:         "Investment",
	BroadcastedEventNameReferredUserPassedKyc: "Community",
	BroadcastedEventNameReceivedAirdrop:       "Community",
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
	BroadcastedEventNameReceivedAirdrop:       "Received Airdrop",
}

// SalesforceTarget represents a target with salesforce api
type SalesforceTarget struct {
	*salesforce.Connector
}

// NewSalesforceTarget constructs a target
func NewSalesforceTarget(sc *salesforce.Connector) *SalesforceTarget {
	return &SalesforceTarget{
		sc,
	}
}

// SendEvent uses salesforce client connector for sending event to analytics
func (st *SalesforceTarget) SendEvent(event *BroadcastedEvent) error {
	_, err := st.Connector.PostUserAction(salesforce.UserActionData{
		Sphere:           eventNameToSphere[event.Name],
		ActionName:       eventNameToActionName[event.Name],
		OccurredAt:       event.Time,
		ActorName:        event.ActorName,
		ActorEmail:       event.ActorEmail,
		InvestmentAmount: event.InvestmentAmount,
		DepositAmount:    event.DepositAmount,
		DepositCurrency:  event.DepositCurrency,
		Referrer:         event.Referral,
		Country:          event.Country,
	})
	if err != nil {
		return errors.Wrap(err, "failed to post event")
	}
	return nil
}
