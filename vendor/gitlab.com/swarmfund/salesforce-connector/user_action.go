package salesforce

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// Event contains both default salesforce and specific to a salesforce account fields defined as columns
type Event struct {
	Name                   string
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

// EventResponse holds data received after SendEvent
type EventResponse struct {
	SalesforceID string   `json:"id"`
	Success      bool     `json:"success"`
	Errors       []string `json:"errors"`
}

var eventsEndpointURL = &url.URL{
	Path: "/services/data/v42.0/sobjects/Website_Action__c/",
}

// PostEvent sends an event to predefined salesforce endpoint, uses now-time if failed to parse timeString
func (c *Connector) PostEvent(sphere string, actionName string, timeString string, actorName string, actorEmail string, investmentAmount int64, depositAmount int64, depositCurrency, referral, country string) (*EventResponse, error) {
	requestStruct := &Event{
		Name:                   "Action",
		PropertyColumn:         "Swarm Invest",
		SphereColumn:           sphere,
		ActionColumn:           actionName,
		ActionDateTimeColumn:   timeString,
		ActorNameColumn:        actorName,
		ActorEmailColumn:       actorEmail,
		InvestmentAmountColumn: investmentAmount,
		DepositAmountColumn:    depositAmount,
		DepositCurrencyColumn:  depositCurrency,
		ReferralColumn:         referral,
		CountryColumn:          country,
	}

	requestBytes, err := json.Marshal(requestStruct)
	if err != nil {
		return nil, err
	}

	statusCode, responseBytes, err := c.client.PostObject(requestBytes, eventsEndpointURL)
	if err != nil {
		return nil, err
	}

	switch statusCode {
	case http.StatusCreated:

		var eventResponse *EventResponse
		err = json.Unmarshal(responseBytes, &eventResponse)
		if err != nil {
			return nil, err
		}

		return eventResponse, nil

	case http.StatusUnauthorized:
		return nil, errors.New("unauthorized")
	case http.StatusBadRequest:
		return nil, errors.New("malformed request sent")
	}

	return nil, errors.New("unknown status code")
}
