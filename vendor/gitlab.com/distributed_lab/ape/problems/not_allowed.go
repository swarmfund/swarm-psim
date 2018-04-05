package problems

import (
	"fmt"
	"net/http"

	"github.com/google/jsonapi"
)

func NotAllowed() *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusUnauthorized),
		Status: fmt.Sprintf("%d", http.StatusUnauthorized),
	}
}
