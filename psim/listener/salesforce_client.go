package listener

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// SalesforceClient is a custom implementation of salesforce client derived from horizon client
type SalesforceClient struct {
	base   *url.URL
	client *http.Client
}

// NewSalesforceClient constructs a SalesforceClient with base url using httpClient
func NewSalesforceClient(client *http.Client, base *url.URL) *SalesforceClient {
	return &SalesforceClient{
		base, client,
	}
}

func (sc *SalesforceClient) resolveURL(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", errors.Wrap(err, "Failed to parse endpoint into URL")
	}

	return sc.base.ResolveReference(u).String(), nil
}

func ignoreError(err error) {
	_ = err
}

func (sc *SalesforceClient) do(request *http.Request, contentType string) (int, []byte, error) {
	request.Header.Set("content-type", contentType)
	request.Header.Set("accept", "application/json")

	response, err := sc.client.Do(request)
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to perform http request")
	}
	if response == nil {
		return 0, nil, errors.New("nil response received")
	}
	defer func() {
		ignoreError(response.Body.Close())
	}()

	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to read response body")
	}

	return response.StatusCode, respBytes, nil
}

// PostURLEncoded sends an x-www-form-urlencoded request using data from req to endpoint
func (sc *SalesforceClient) PostURLEncoded(endpoint string, req []byte) (statusCode int, response []byte, err error) {
	endpoint, err = sc.resolveURL(endpoint)
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to resolve url")
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewReader(req))
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to create POST http.Request")
	}

	statusCode, responseBytes, err := sc.do(request, "application/x-www-form-urlencoded")
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to do the request")
	}

	return statusCode, responseBytes, nil
}

// PostJSON sends json data to endpoint
func (sc *SalesforceClient) PostJSON(endpoint string, req interface{}, accessToken string) (statusCode int, response []byte, err error) {
	requestBytes, err := json.Marshal(req)
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to marshal request into JSON bytes")
	}

	endpoint, err = sc.resolveURL(endpoint)
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to resolve url")
	}

	request, err := http.NewRequest("POST", endpoint, bytes.NewReader(requestBytes))
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to create POST http.Request")
	}
	if accessToken == "" {
		return 0, nil, errors.New("got empty access token")
	}
	request.Header.Set("authorization", "Bearer "+accessToken)

	statusCode, responseBB, err := sc.do(request, "application/json")
	if err != nil {
		return 0, nil, errors.Wrap(err, "Failed to do the request")
	}

	return statusCode, responseBB, nil
}

const authEndpoint = "services/oauth2/token"

const clientID = "3MVG9RHx1QGZ7OsjxpWOBJ4UsJLqMgjd8FcZ8K9fCEY0YhJ5Av_FbGhmSQLfWoaDj3AQO9rcxwvbQNLDF_LH5"

// PostAuthRequest sends auth request to salesforce api using username, password and clientSecret from arguments and hardcoded clientId
func (sc *SalesforceClient) PostAuthRequest(username string, password string, clientSecret string) (statusCode int, response []byte, err error) {
	endpointString := sc.base.String() + authEndpoint
	requestString := "username=" + username + "&client_secret=" + clientSecret + "&password=" + password + "&grant_type=password&client_id=" + clientID
	return sc.PostURLEncoded(endpointString, []byte(requestString))
}

// SalesforceAuthResponse contains auth response data, only AccessToken is useful
type SalesforceAuthResponse struct {
	AccessToken string `json:"access_token"`
}

// SalesforceEvent contains both default salesforce and specific to a salesforce account fields defined as columns
type SalesforceEvent struct {
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

const eventsEndpoint = "/services/data/v42.0/sobjects/Website_Action__c/"

// PostEvent sends an event to predefined salesforce endpoint, uses now-time if failed to parse timeString
func (sc *SalesforceClient) PostEvent(accessToken string, sphere string, actionName string, timeString string, actorName string, actorEmail string, investmentAmount int64, investmentCountry string) (statusCode int, response []byte, err error) {
	endpointString := sc.base.String() + eventsEndpoint
	requestStruct := &SalesforceEvent{
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
	return sc.PostJSON(endpointString, requestStruct, accessToken)
}

// SalesforceEventResponse holds data received after SendEvent
type SalesforceEventResponse struct {
	SalesforceID string   `json:"id"`
	Success      bool     `json:"success"`
	Errors       []string `json:"errors"`
}
