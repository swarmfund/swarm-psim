package responses

import "gitlab.com/swarmfund/horizon-connector/v2/internal/resources"

type Users struct {
	Data   []resources.User `json:"data"`
	Links  Links            `json:"links"`
	Errors []Error          `json:"errors"`
}
