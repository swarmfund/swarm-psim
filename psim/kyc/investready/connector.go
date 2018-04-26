package investready

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"io/ioutil"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	pageSize = 100 // 100 is max
)

var (
	errUnauthorised = errors.New("Got Unauthorised status code.")
)

type Connector struct {
	log *logan.Entry
	config ConnectorConfig

	client      *http.Client
	accessToken string
}

func NewConnector(log *logan.Entry, config ConnectorConfig) *Connector {
	return &Connector{
		log: log,
		config: config,

		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// TODO Comment
func (c *Connector) ObtainUserToken(oauthCode string) (userAccessToken string, err error) {
	req := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
		Code         string `json:"code"`
	}{
		GrantType:    "authorization_code",
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		// TODO
		RedirectURI: "https://invest.swarm.fund",
		Code:        oauthCode,
	}

	reqBB, err := json.Marshal(req)
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal request")
	}

	resp, err := c.client.Post(fmt.Sprintf("%s/oauth/token", c.config.URL), "application/json", bytes.NewReader(reqBB))
	if err != nil {
		return "", errors.Wrap(err, "Failed to send http POST")
	}
	fields := logan.F{
		"response_status_code": resp.StatusCode,
	}

	respBB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read response body into bytes", fields)
	}
	fields["raw_response"] = string(respBB)

	if resp.StatusCode != http.StatusOK {
		return "", errors.From(errors.New("Received response with unsuccessful status code."), fields)
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(respBB, &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed unmarshal response bytes", fields)
	}

	return response.AccessToken, nil
}

// TODO Comment
// UserHash can possibly return "",nil if the response won't contain UserHash in proper field.
func (c *Connector) UserHash(userAccessToken string) (userHash string, err error) {
	req := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: userAccessToken,
	}

	reqBB, err := json.Marshal(req)
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal request")
	}

	resp, err := c.client.Post(fmt.Sprintf("%s/api/me.json", c.config.URL), "application/json", bytes.NewReader(reqBB))
	if err != nil {
		return "", errors.Wrap(err, "Failed to send http POST")
	}
	fields := logan.F{
		"response_status_code": resp.StatusCode,
	}

	respBB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read response body into bytes", fields)
	}
	fields["raw_response"] = string(respBB)

	if resp.StatusCode != http.StatusOK {
		return "", errors.From(errors.New("Received response with unsuccessful status code."), fields)
	}

	var response struct {
		Data struct {
			Person struct {
				Hash string `json:"hash"`
			} `json:"person"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBB, &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal response bytes", fields)
	}

	return response.Data.Person.Hash, nil
}

func (c *Connector) ListAllSyncedUsers() ([]User, error) {
	var page = 0
	allUsers := make([]User, 0)

	for {
		usersPage, err := c.getSyncedUsersPageWithTokenRefresh(page)
		if err != nil {
			return allUsers, errors.Wrap(err, "Failed to get Users page", logan.F{
				"page": page,
			})
		}

		allUsers = append(allUsers, usersPage...)
		page += 1

		if len(usersPage) < pageSize {
			// Last page obtained.
			return allUsers, nil
		}
	}
}

func (c *Connector) getSyncedUsersPageWithTokenRefresh(page int) ([]User, error) {
	users, err := c.getSyncedUsersPage(page)
	if err == errUnauthorised {
		c.log.WithError(err).Info("Got Unauthorised status code - refreshing Client AccessToken.")

		err = c.refreshClientToken()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to refresh Client AccessToken")
		}

		return c.getSyncedUsersPage(page)
	}

	return users, err
}

func (c *Connector) getSyncedUsersPage(page int) ([]User, error) {
	req := struct {
		Page  int `json:"page"`
		Count int `json:"count"`
	}{
		Page:  page,
		Count: pageSize,
	}

	reqBB, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal request")
	}

	resp, err := c.client.Post(fmt.Sprintf("%s/api/me.json", c.config.URL), "application/json", bytes.NewReader(reqBB))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send http POST")
	}
	fields := logan.F{
		"response_status_code": resp.StatusCode,
	}

	respBB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read response body into bytes", fields)
	}
	fields["raw_response"] = string(respBB)

	if resp.StatusCode == http.StatusUnauthorized {
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.From(errors.New("Received response with unsuccessful status code."), fields)
	}

	var response struct {
		Data struct {
			// Response also contains `pagination` object with pages, current_page, limit and user_count
			Users []User `json:"users"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBB, &response)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal response bytes", fields)
	}

	return response.Data.Users, nil
}

func (c *Connector) refreshClientToken() error {
	newAccessToken, err := c.getClientToken()
	if err != nil {
		return errors.Wrap(err, "Failed to get Client token")
	}

	if newAccessToken == "" {
		return errors.New("Obtained empty Client token.")
	}

	c.accessToken = newAccessToken
	return nil
}

func (c *Connector) getClientToken() (clientAccessToken string, err error) {
	req := struct {
		GrantType    string `json:"grant_type"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
	}{
		GrantType:    "client_credentials",
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		// TODO
		RedirectURI: "https://invest.swarm.fund",
	}

	reqBB, err := json.Marshal(req)
	if err != nil {
		return "", errors.Wrap(err, "Failed to marshal request")
	}

	resp, err := c.client.Post(fmt.Sprintf("%s/oauth/token", c.config.URL), "application/json", bytes.NewReader(reqBB))
	if err != nil {
		return "", errors.Wrap(err, "Failed to send http POST")
	}
	fields := logan.F{
		"response_status_code": resp.StatusCode,
	}

	respBB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read response body into bytes", fields)
	}
	fields["raw_response"] = string(respBB)

	if resp.StatusCode != http.StatusOK {
		return "", errors.From(errors.New("Received response with unsuccessful status code."), fields)
	}

	var response struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(respBB, &response)
	if err != nil {
		return "", errors.Wrap(err, "Failed to unmarshal response bytes", fields)
	}

	return response.AccessToken, nil
}
