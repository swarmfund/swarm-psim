package handlers

import (
	"net/http/httptest"
	"testing"

	"net/http"
	"strings"
)

func TestStripeChargeHandler(t *testing.T) {
	t.Run("empty body", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", nil)
		req.Header.Set("content-type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(StripeChargeHandler)

		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf(
				"expected %v status code got %v", http.StatusBadRequest, status)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{{`))
		req.Header.Set("content-type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(StripeChargeHandler)

		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf(
				"expected %v status code got %v", http.StatusBadRequest, status)
		}
	})
}
