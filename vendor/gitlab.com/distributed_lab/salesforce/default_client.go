package salesforce

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// client is a custom salesforce client implementation
type client struct {
	httpClient  *http.Client
	apiURL      *url.URL
	secret      string
	id          string
	accessToken string
	username    string
	password    string
}

func (c *client) authenticate() error {
	authResponse, err := c.postAuthRequest(c.username, c.password)
	if err != nil {
		return errors.Wrap(err, "failed to authenticate")
	}
	if authResponse == nil {
		return errors.New("got empty auth response")
	}

	c.accessToken = authResponse.AccessToken
	return nil
}

// TODO: make configurable
var authEndpointURL = &url.URL{
	Path: "services/oauth2/token",
}

// AuthResponse contains auth response data, only AccessToken is useful
type authResponse struct {
	AccessToken string `json:"access_token"`
}

// Custom errors
var (
	errMalformedRequest = errors.New("malformed request sent")
	errInternal         = errors.New("something bad happened")
)

// PostAuthRequest sends auth request to salesforce api using username, password and clientSecret from arguments
// the access token is not expected to be expired
func (c *client) postAuthRequest(username string, password string) (*authResponse, error) {
	requestString := url.Values{}
	requestString.Set("username", username)
	requestString.Set("client_secret", c.secret)
	requestString.Set("password", password)
	requestString.Set("grant_type", "password")
	requestString.Set("client_id", c.id)

	authURL := c.apiURL.ResolveReference(authEndpointURL)

	response, err := c.httpClient.Post(authURL.String(), "application/x-www-form-urlencoded", strings.NewReader(requestString.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to send auth post request")
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read auth response body")
	}

	switch response.StatusCode {
	case http.StatusOK:
		var authResponse authResponse
		err = json.Unmarshal(responseBytes, &authResponse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal auth response json", logan.F{
				"response": string(responseBytes),
			})
		}
		return &authResponse, nil
	case http.StatusBadRequest:
		return nil, errors.From(errMalformedRequest, logan.F{"response_body": string(responseBytes)})
	default:
		return nil, errors.From(errInternal, logan.F{
			"response_body": string(responseBytes),
			"status_code":   response.StatusCode,
		})
	}
}

// PostObject sends json data to an endpoint
func (c *client) PostObject(json []byte, endpoint *url.URL) (statusCode int, body []byte, err error) {
	endpointURL := c.apiURL.ResolveReference(endpoint)

	req, err := http.NewRequest("POST", endpointURL.String(), bytes.NewReader(json))
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	if c.accessToken == "" {
		err := c.authenticate()
		if err != nil {
			return 0, nil, errors.Wrap(err, "failed to authenticate by username and password")
		}
	}
	req.Header.Set("authorization", "Bearer "+c.accessToken)

	response, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to perform request")
	}

	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, nil, errors.Wrap(err, "failed to read response body")
	}

	return response.StatusCode, responseBytes, nil
}
