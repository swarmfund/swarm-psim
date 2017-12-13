package main

import (
	"net/url"

	. "gitlab.com/distributed_lab/notificator-server/client"
)

func main() {
	client := NewConnector(
		Pair{Secret: "secretkey", Public: "publickey"},
		url.URL{Scheme: "http", Host: "localhost:9009"})
	payload := EmailRequestPayload{
		Subject:     "subject",
		Destination: "foo@bar.egg",
		Message:     "ohai",
	}
	response, err := client.Send(1, "user_token", payload)

	if err != nil {
		if err == ErrInternal {
			// you are doing something wrong
		}
		// request was not send, probably due to transport layer issues
		// retry logic is up to the user
	}

	if !response.IsSuccess() {
		// something wrong on application level
		if response.IsPermanent() {
			if !response.Authenticated() {
				// authenticated failed, either key or signature is not valid
			}
			// issue is permanent so there is no reason to retry request
			// probably request was already submitted or malformed in some way
			// just deal with it
		}

		if retryIn := response.RetryIn(); retryIn != nil {
			// your request was rate limited or for some reason rejected
			// `retryIn` contains suggested `time.Duration` for next attempt
		} else {
			// you could retry but service didn't provide any hints
		}
	}
}
