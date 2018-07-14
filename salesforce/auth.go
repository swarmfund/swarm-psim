package salesforce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/distributed_lab/logan"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// AuthResponse contains auth response data, only AccessToken is useful
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// EmptyAuthResponse is used for signaling about special conditions
var (
	EmptyAuthResponse   = AuthResponse{}
	ErrMalformedRequest = errors.New("malformed request sent")
	ErrInternal         = errors.New("something bad happened")
)

// PostAuthRequest sends auth request to salesforce api using username, password and clientSecret from arguments and hardcoded clientId
func (c *Client) PostAuthRequest(username string, password string) (AuthResponse, error) {
	requestString := url.Values{}
	requestString.Set("username", username)
	requestString.Set("client_secret", c.secret)
	requestString.Set("password", password)
	requestString.Set("grant_type", "password")
	requestString.Set("client_id", c.id)

	authURL := c.apiURL.ResolveReference(authEndpointURL)

	response, err := c.httpClient.Post(authURL.String(), "application/x-www-form-urlencoded", strings.NewReader(requestString.Encode()))
	if err != nil {
		return EmptyAuthResponse, errors.Wrap(err, "failed to post")
	}

	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return EmptyAuthResponse, errors.Wrap(err, "failed to read auth request body")
	}

	switch response.StatusCode {
	case http.StatusOK:
		var authResponse AuthResponse
		err = json.Unmarshal(responseBytes, &authResponse)
		if err != nil {
			return EmptyAuthResponse, errors.Wrap(err, "failed to unmarshal auth response json")
		}
		return authResponse, nil
	case http.StatusBadRequest:
		return EmptyAuthResponse, errors.From(ErrMalformedRequest, logan.F{"response_body": string(responseBytes)})
	default:
		return EmptyAuthResponse, errors.From(ErrInternal, logan.F{
			"response_body": string(responseBytes),
			"status_code":   response.StatusCode,
		})
	}
}
