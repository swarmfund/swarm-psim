package salesforce

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"
)

type Report struct {
	PositiveColumn  int    `json:"Total_Individual_Accounts_Positive_Funds__c"`
	ZeroColumn      int    `json:"Total_Individual_Accounts_Zero_Funds__c"`
	SWMAmountColumn int64  `json:"SWM_Amount__c"`
	ThresholdColumn int64  `json:"X_Value__c"`
	AboveColumn     int    `json:"Total_Values_Above_X__c"`
	BelowColumn     int    `json:"Total_Values_Below_X__c"`
	DateColumn      string `json:"Date__c"` // time.Format("2006-01-02")
}

type ReportResponse struct {
	SalesforceID string   `json:"id"`
	Success      bool     `json:"success"`
	Errors       []string `json:"errors"`
}

var reportEndpointURL = &url.URL{
	Path: "/services/data/v42.0/sobjects/Swarm_Statistic__c/",
}

func (c *Connector) PostReport(positive, zero int, swmAmount, threshold int64, above, below int, date *time.Time) (*EventResponse, error) {
	requestStruct := &Report{
		PositiveColumn:  positive,
		ZeroColumn:      zero,
		SWMAmountColumn: swmAmount,
		ThresholdColumn: threshold,
		AboveColumn:     above,
		BelowColumn:     below,
		DateColumn:      date.Format("2006-01-02"),
	}

	requestBytes, err := json.Marshal(requestStruct)
	if err != nil {
		return nil, err
	}

	statusCode, responseBytes, err := c.client.PostObject(requestBytes, reportEndpointURL)
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
