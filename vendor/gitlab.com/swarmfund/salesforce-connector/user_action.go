package salesforce

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const salesforceTimeLayout = "2006-01-02T15:04:05.999-0700"

// UserAction contains both default salesforce and specific to a salesforce account fields defined as columns
type userAction struct {
	Name                   string `json:"name"`
	PropertyColumn         string `json:"Property__c"`
	SphereColumn           string `json:"Sphere__c"`
	ActionColumn           string `json:"Action__c"`
	ActionDateTimeColumn   string `json:"Action_Date_Time__c"` // time.Format("2006-01-02T15:04:05.999-0700")
	ActorNameColumn        string `json:"Actor_Name__c"`
	ActorEmailColumn       string `json:"Actor_Email__c"`
	InvestmentAmountColumn int64  `json:"Investment_Amount__c"`
	DepositAmountColumn    int64  `json:"Deposit_Amount__c"`
	DepositCurrencyColumn  string `json:"Deposit_Currency__c"`
	ReferralColumn         string `json:"Referral__c"`
	CountryColumn          string `json:"Country__c"`
}

var userActionsEndpointURL = &url.URL{
	Path: "/services/data/v42.0/sobjects/Website_Action__c/",
}

type UserActionData struct {
	Sphere           string
	ActionName       string
	OccurredAt       time.Time
	ActorName        string
	ActorEmail       string
	InvestmentAmount int64
	DepositAmount    int64
	DepositCurrency  string
	Referrer         string
	Country          string
}

// PostUserAction sends an event to predefined salesforce endpoint, uses now-time if failed to parse timeString
func (c *Connector) PostUserAction(data UserActionData) (*PostObjectResponse, error) {
	requestStruct := &userAction{
		Name:                   "Action",
		PropertyColumn:         "Swarm Invest",
		SphereColumn:           data.Sphere,
		ActionColumn:           data.ActionName,
		ActionDateTimeColumn:   data.OccurredAt.Format(salesforceTimeLayout),
		ActorNameColumn:        data.ActorName,
		ActorEmailColumn:       data.ActorEmail,
		InvestmentAmountColumn: data.InvestmentAmount,
		DepositAmountColumn:    data.DepositAmount,
		DepositCurrencyColumn:  data.DepositCurrency,
		ReferralColumn:         data.Referrer,
		CountryColumn:          data.Country,
	}

	requestBytes, err := json.Marshal(requestStruct)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request struct", logan.F{
			"request": requestStruct,
		})
	}

	statusCode, responseBytes, err := c.client.PostObject(requestBytes, userActionsEndpointURL)
	if err != nil {
		return nil, errors.Wrap(err, "faield to post report object", logan.F{
			"request": string(requestBytes),
		})
	}

	switch statusCode {
	case http.StatusCreated:
		var eventResponse *PostObjectResponse
		err = json.Unmarshal(responseBytes, &eventResponse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal response", logan.F{
				"response": string(responseBytes),
			})
		}
		return eventResponse, nil
	case http.StatusUnauthorized:
		return nil, errors.New("unauthorized")
	case http.StatusBadRequest:
		return nil, errors.From(errors.New("malformed request sent"), logan.F{
			"response": string(responseBytes),
			"request":  string(requestBytes),
		})
	default:
		return nil, errors.From(errors.New("unknown status code"), logan.F{
			"status_code": statusCode,
			"response":    string(responseBytes),
		})
	}
}
