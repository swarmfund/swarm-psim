package handlers

import (
	"net/http"

	"github.com/stripe/stripe-go/client"
	"gitlab.com/distributed_lab/logan/v3"
	"context"
)

type ctxKey uint8

const (
	stripeCtxKey ctxKey = iota
	logCtxKey
)

func Stripe(r *http.Request) *client.API {
	return r.Context().Value(stripeCtxKey).(*client.API)
}

func PutStripe(r *http.Request, stripe *client.API) {
	ctx := context.WithValue(r.Context(), stripeCtxKey, stripe)
	r.WithContext(ctx)
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func PutLog(r *http.Request, log *logan.Entry) {
	ctx := context.WithValue(r.Context(), logCtxKey, log)
	r.WithContext(ctx)
}

// TODO Add Horizon connector context methods
