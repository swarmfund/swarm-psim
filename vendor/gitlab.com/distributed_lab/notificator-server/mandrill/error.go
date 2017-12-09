package mandrill

import (
	"fmt"
)

type Error struct {
	Status string `json:"status"`
	Code int `json:"code"`
	Name string `json:"name"`
	Message string `json:"message"`
}

func (e *Error) ToError() error {
	return fmt.Errorf("Status: %s; Name: %s; Message: %s", e.Status, e.Name, e.Message)
}
