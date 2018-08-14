package salesforce

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type BalanceReportData struct {
	Positive  int
	Zero      int
	SWMAmount int64
	Threshold int64
	Above     int
	Below     int
	Date      time.Time
}

type report struct {
	PositiveColumn  int    `json:"Total_Individual_Accounts_Positive_Funds__c"`
	ZeroColumn      int    `json:"Total_Individual_Accounts_Zero_Funds__c"`
	SWMAmountColumn int64  `json:"SWM_Amount__c"`
	ThresholdColumn int64  `json:"X_Value__c"`
	AboveColumn     int    `json:"Total_Values_Above_X__c"`
	BelowColumn     int    `json:"Total_Values_Below_X__c"`
	DateColumn      string `json:"Date__c"` // time.Format("2006-01-02")
}

var reportEndpointURL = &url.URL{
	Path: "/services/data/v42.0/sobjects/Swarm_Statistic__c/",
}

func (c *Connector) PostReport(reportData BalanceReportData) (*PostObjectResponse, error) {
	requestStruct := report{
		PositiveColumn:  reportData.Positive,
		ZeroColumn:      reportData.Zero,
		SWMAmountColumn: reportData.SWMAmount,
		ThresholdColumn: reportData.Threshold,
		AboveColumn:     reportData.Above,
		BelowColumn:     reportData.Below,
		DateColumn:      reportData.Date.Format("2006-01-02"),
	}
	requestBytes, err := json.Marshal(requestStruct)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal request struct to bytes while posting report", logan.F{
			"request": requestStruct,
		})
	}

	statusCode, responseBytes, err := c.client.PostObject(requestBytes, reportEndpointURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to post report object", logan.F{
			"request": string(requestBytes),
		})
	}

	switch statusCode {
	case http.StatusCreated:
		var eventResponse *PostObjectResponse
		err = json.Unmarshal(responseBytes, &eventResponse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal report post response", logan.F{
				"response": string(responseBytes),
			})
		}
		return eventResponse, nil
	case http.StatusUnauthorized:
		return nil, errors.New("unauthorized")
	case http.StatusBadRequest:
		return nil, errors.From(errors.New("malformed request sent"), logan.F{
			"request":  string(requestBytes),
			"response": string(responseBytes),
		})
	default:
		return nil, errors.From(errors.New("unknown status code"), logan.F{
			"status_code": statusCode,
			"response":    string(responseBytes),
		})
	}
}
