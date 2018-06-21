package salesforce

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Event contains both default salesforce and specific to a salesforce account fields defined as columns
type Event struct {
	Name                    string
	PropertyColumn          string `json:"Property__c"`
	SphereColumn            string `json:"Sphere__c"`
	ActionColumn            string `json:"Action__c"`
	ActionDateTimeColumn    string `json:"Action_Date_Time__c"` // time.Format("2006-01-02T15:04:05.999-0700")
	ActorNameColumn         string `json:"Actor_Name__c"`
	ActorEmailColumn        string `json:"Actor_Email__c"`
	InvestmentAmountColumn  int64  `json:"Investment_Amount__c"`
	InvestmentCountryColumn string `json:"Investment_Country__c"`
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
func (c *Client) PostEvent(sphere string, actionName string, timeString string, actorName string, actorEmail string, investmentAmount int64, investmentCountry string) (*EventResponse, error) {
	endpointURL := c.apiURL.ResolveReference(eventsEndpointURL)
	requestStruct := &Event{
		Name:                    "Action",
		PropertyColumn:          "Swarm Invest",
		SphereColumn:            sphere,
		ActionColumn:            actionName,
		ActionDateTimeColumn:    timeString,
		ActorNameColumn:         actorName,
		ActorEmailColumn:        actorEmail,
		InvestmentAmountColumn:  investmentAmount,
		InvestmentCountryColumn: investmentCountry,
	}

	requestBytes, err := json.Marshal(requestStruct)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpointURL.String(), bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", "Bearer "+c.accessToken)
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	// TODO auth.go
	switch response.StatusCode {
	case http.StatusCreated:
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

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
