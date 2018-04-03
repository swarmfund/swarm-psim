package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/go/signcontrol"
)

func RenderDoormanErr(w http.ResponseWriter, err error) {
	switch err {
	case signcontrol.ErrNotSigned, signcontrol.ErrValidUntil, signcontrol.ErrExpired, signcontrol.ErrSignerKey, signcontrol.ErrSignature, signcontrol.ErrNotAllowed:
		ape.RenderErr(w, problems.NotAllowed())
	case nil:
		panic("expected not nil error")
	default:
		panic(errors.Wrap(err, "unexpected doorman error"))
	}
}
