package utils

import "testing"

func TestParseBalanceID(t *testing.T) {
	_, err := ParseBalanceID("yoba")
	if err == nil {
		t.Error("should return err")
	}
	balanceID, err := ParseBalanceID("BDRYPVZ63SR7V2G46GKRGABJD3XPDNWQ4B4PQPJBTTDUEAKH5ZECP4UG")
	if err != nil {
		t.Error("error should be nil")
	}
	if balanceID.AsString() != "BDRYPVZ63SR7V2G46GKRGABJD3XPDNWQ4B4PQPJBTTDUEAKH5ZECP4UG" {
		t.Error("value is wrong")
	}
}
