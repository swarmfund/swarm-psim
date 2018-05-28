package listener

import (
	"time"

	"gitlab.com/swarmfund/psim/psim/listener/internal"
)

type BroadcastedEvent = internal.BroadcastedEvent
type BroadcastedEventName = internal.BroadcastedEventName
type Processor = internal.Processor
type TxData = internal.TxData

// TODO rename Output to Broadcasted
const (
	OutputEventNameKycCreated            BroadcastedEventName = "kyc_created"
	OutputEventNameKycUpdated            BroadcastedEventName = "kyc_updated"
	OutputEventNameKycRejected           BroadcastedEventName = "kyc_rejected"
	OutputEventNameKycApproved           BroadcastedEventName = "kyc_approved"
	OutputEventNameUserReferred          BroadcastedEventName = "user_referred"
	OutputEventNameFundsWithdrawn        BroadcastedEventName = "funds_withdrawn"
	OutputEventNamePaymentV2Received     BroadcastedEventName = "payment_v2_received"
	OutputEventNamePaymentV2Sent         BroadcastedEventName = "payment_v2_sent"
	OutputEventNamePaymentReceived       BroadcastedEventName = "payment_received"
	OutputEventNamePaymentSent           BroadcastedEventName = "payment_sent"
	OutputEventNameFundsDeposited        BroadcastedEventName = "funds_deposited"
	OutputEventNameFundsInvested         BroadcastedEventName = "funds_invested"
	OutputEventNameReferredUserPassedKyc BroadcastedEventName = "referred_user_passed_kyc"
)

const (
	defaultServiceRetryTimeIncrement = 1 * time.Second
	defaultMaxServiceRetryTime       = 30 * time.Second
)

const defaultMixpanelURL = ""
