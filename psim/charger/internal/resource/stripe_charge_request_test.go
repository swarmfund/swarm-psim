package resource

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStripeChargeRequest_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name string
		data string
		err  error
		res  StripeChargeRequest
	}{
		{"empty", `{}`, nil, StripeChargeRequest{}},
		{"zero amount", `{"amount": "0"}`, nil, StripeChargeRequest{}},
		{"positive amount", `{"amount": "1.0"}`, nil, StripeChargeRequest{Amount: 10000}},
	}

	for _, tc := range testCases {
		var r StripeChargeRequest
		t.Run(tc.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tc.data), &r)
			if err != tc.err {
				t.Errorf("got %v expected %v", err, tc.err)
				return
			}
			if !reflect.DeepEqual(tc.res, r) {
				t.Errorf("got %v expected %v", r, tc.res)
				return
			}
		})
	}
}

func TestStripeChargeRequest_Validate(t *testing.T) {
	testCases := []struct {
		name string
		obj  StripeChargeRequest
		err  bool
	}{
		{"empty", StripeChargeRequest{}, true},
		{"valid", StripeChargeRequest{"token", minStripeAmountValue * 2, "reference", "", ""}, false},
		{"low amount", StripeChargeRequest{"token", 1, "reference", "", ""}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.obj.Validate()
			if !tc.err && err != nil {
				t.Errorf("got %v expected nil", err)
				return
			}
			if tc.err && err == nil {
				t.Error("got nil expected error")
				return
			}
		})
	}
}
