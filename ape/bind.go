package ape

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Validator interface {
	Validate() error
}

var (
	ErrNilBody = errors.New("nil body")
)

func Bind(r *http.Request, val interface{}) error {
	/*ct := r.Header.Get("content-type")
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return problems.UnsupportedMediaType(err)
	}

	if mt != "application/json" {
		return problems.UnsupportedMediaType(nil)
	}*/

	if r.ContentLength == 0 {
		return ErrNilBody
	}

	err := json.NewDecoder(r.Body).Decode(val)
	if err != nil {
		return err
	}

	if validator, ok := val.(Validator); ok {
		err = validator.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
