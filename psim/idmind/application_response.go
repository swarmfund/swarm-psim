package idmind

import (
	"strconv"
	"strings"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// TODO GetLoganFields implementations

// Don't ask why json keys are so strange - that's how IdentityMind works.
type ApplicationResponse struct {
	CheckApplicationResponse

	// The most interesting fields
	PolicyResult FraudResult `json:"res"`

	EDNAPolicyResult Reputation   `json:"user"`

	FraudResult               FraudResult `json:"frp"`
	FiredFraudRule            string      `json:"frn"`
	FiredFraudRuleDescription string      `json:"frd"`

	ReputationReason string `json:"erd"`
}

func (r ApplicationResponse) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"m_tx_id":               r.MTxID,
		"tx_id":                 r.TxID,
		"kyc_state":             r.KYCState,
		"policy_result":         r.PolicyResult,
		"edna_policy_result":    r.EDNAPolicyResult,
		"edna_score_card":       r.EDNAScoreCard,
		"fraud_result":          r.FraudResult,
		"fired_fraud_rule":      r.FiredFraudRule,
		"fired_fraud_rule_desc": r.FiredFraudRuleDescription,
		"reputation_reason":     r.ReputationReason,
		"result_code":           r.ResultCodes,
	}
}

type CheckApplicationResponse struct {
	// I still donno what's the difference between MTxID and TxID.
	MTxID string `json:"mtid"`
	TxID  string `json:"tid"`

	// The most interesting fields
	KYCState KYCState `json:"state"`

	EDNAScoreCard    ExtScoreCard `json:"ednaScoreCard"`
	ResultCodes string `json:"rcd"`
}

func (r CheckApplicationResponse) GetResultCodes() ([]int, error) {
	ss := strings.Split(r.ResultCodes, ",")

	var result []int
	for i, s := range ss {
		number, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to parse integer from the result code", logan.F{
				"result_code_i":  i,
				"raw_code_value": s,
			})
		}

		result = append(result, number)
	}

	return result, nil
}

func (r CheckApplicationResponse) GetLoganFields() map[string]interface{} {
	return map[string]interface{}{
		"m_tx_id":               r.MTxID,
		"tx_id":                 r.TxID,
		"kyc_state":             r.KYCState,
		"edna_score_card":       r.EDNAScoreCard,
		"result_code":           r.ResultCodes,
	}
}

type Reputation string

const (
	TrustedReputation    Reputation = "TRUSTED"
	UnknownReputation    Reputation = "UNKNOWN"
	SuspiciousReputation Reputation = "SUSPICIOUS"
	BadReputation        Reputation = "BAD"
)

type FraudResult string

const (
	AccepFraudResult        = "ACCEPT"
	ManualReviewFraudResult = "MANUAL_REVIEW"
	DenyFraudResult         = "DENY"
)

type KYCState string

const (
	AcceptedKYCState    KYCState = "A"
	UnderReviewKYCState KYCState = "R"
	RejectedKYCState    KYCState = "D"
)

// TODO Implement GetLoganFields somehow
type ExtScoreCard struct {
	TestResults           []ConditionResult     `json:"sc"`
	EvaluatedTestResults  []ConditionResult     `json:"etr"`
	FraudPolicyEvaluation ExtEvalResult         `json:"er"`
	AutomatedReview       AutomatedReviewResult `json:"ar"`
}

// TODO Implement GetLoganFields
type ConditionResult struct {
	Test           string `json:"test"`
	Details        string `json:"details"`
	Fired          bool   `json:"fired"`
	Timestamp      uint64 `json:"ts"`
	KYCStage       string `json:"stage"`
	WaitingForData bool   `json:"waitingForData"`
	//Condition Condition `json:"condition"` // Present in response, but not described in docs.
}

// TODO Implement GetLoganFields
type ExtEvalResult struct {
	ReportedRule ExtRule `json:"reportedRule"`
	Profile      string  `json:"profile"`
}

// TODO Implement GetLoganFields
type ExtRule struct {
	Name                 string            `json:"name"`
	Description          string            `json:"description"`
	Details              string            `json:"details"`
	ResultCode           FraudResult       `json:"resultCode"`
	RuleAssertionResults []ConditionResult `json:"testResults"`
	RuleId               int               `json:"ruleId"`
}

// TODO Implement GetLoganFields
type AutomatedReviewResult struct {
	Result          ReviewResult `json:"result"`
	RuleID          string       `json:"ruleId"`
	RuleName        string       `json:"ruleName"`
	RuleDescription string       `json:"ruleDescription"`
}

type ReviewResult string

const (
	ErrorReviewResult         ReviewResult = "ERROR"
	NoPolicyReviewResult      ReviewResult = "NO_POLICY"
	DisabledReviewResult      ReviewResult = "DISABLED"
	FilteredReviewResult      ReviewResult = "FILTERED"
	PendingReviewResult       ReviewResult = "PENDING"
	FailReviewResult          ReviewResult = "FAIL"
	IndeterminateReviewResult ReviewResult = "INDETERMINATE"
	SuccessReviewResult       ReviewResult = "SUCCESS"
)
