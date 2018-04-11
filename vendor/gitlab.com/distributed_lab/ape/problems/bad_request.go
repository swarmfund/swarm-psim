package problems

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/google/jsonapi"
	"github.com/pkg/errors"
)

func BadRequest(err error) []*jsonapi.ErrorObject {
	errs := []*jsonapi.ErrorObject{}
	switch reason := errors.Cause(err); reason {
	case io.EOF:
		errs = append(errs, &jsonapi.ErrorObject{
			Title:  http.StatusText(http.StatusBadRequest),
			Status: fmt.Sprintf("%d", http.StatusBadRequest),
			Detail: "Request body were expected",
		})
	default:
		switch terr := reason.(type) {
		case validation.Errors:
			for key, value := range terr {
				errs = append(errs, &jsonapi.ErrorObject{
					Title:  http.StatusText(http.StatusBadRequest),
					Status: fmt.Sprintf("%d", http.StatusBadRequest),
					Meta: &map[string]interface{}{
						"field": key,
						"error": value.Error(),
					},
				})
			}
		default:
			errs = append(errs, &jsonapi.ErrorObject{
				Title:  http.StatusText(http.StatusBadRequest),
				Status: fmt.Sprintf("%d", http.StatusBadRequest),
				Detail: "Your request was invalid in some way",
			})
		}
	}
	return errs
}
