package salesforce

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// AuthResponse contains auth response data, only AccessToken is useful
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// EmptyAuthResponse is used for signaling about special conditions
var EmptyAuthResponse = AuthResponse{}

// PostAuthRequest sends auth request to salesforce api using username, password and clientSecret from arguments and hardcoded clientId
func (c *Client) PostAuthRequest() (AuthResponse, error) {
	requestString := url.Values{}
	requestString.Set("username", c.username)
	requestString.Set("client_secret", c.secret)
	requestString.Set("password", c.password)
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
		return EmptyAuthResponse, err
	}

	var authResponse AuthResponse
	err = json.Unmarshal(responseBytes, authResponse)
	if err != nil {
		return EmptyAuthResponse, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		return authResponse, nil
	case http.StatusBadRequest:
		return EmptyAuthResponse, errors.New("malformed request sent")
	default:
		return EmptyAuthResponse, nil
	}
}
