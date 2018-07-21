package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/tokend/go/keypair"
	"gitlab.com/tokend/go/signcontrol"
)

type TestClient struct {
	t      *testing.T
	ts     *httptest.Server
	signer keypair.KP
}

func (c *TestClient) RandomSigner() *TestClient {
	c.signer, _ = keypair.Random()
	return c
}

func (c *TestClient) Signer(signer keypair.KP) *TestClient {
	c.signer = signer
	return c
}

func Client(t *testing.T, ts *httptest.Server) *TestClient {
	return &TestClient{
		t:  t,
		ts: ts,
	}

}

func (c *TestClient) Do(method, path, body string) *http.Response {
	c.t.Helper()
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.ts.URL, path), bytes.NewReader([]byte(body)))
	if err != nil {
		c.t.Fatal(err)
	}

	if c.signer != nil {
		if err := signcontrol.SignRequest(request, c.signer); err != nil {
			c.t.Fatal(err)
		}
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		c.t.Fatal(err)
	}
	return response
}
