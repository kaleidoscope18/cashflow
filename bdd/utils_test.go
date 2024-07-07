package bdd

import (
	"cashflow/models"
	"testing"
)

func TestUnmarshalWrappedBody(t *testing.T) {
	var result models.Balance

	// byte array version of {"data":{"createBalance":{"date":"2022/10/15","amount":1000}}}
	data := []byte{123, 34, 100, 97, 116, 97, 34, 58, 123, 34, 99, 114, 101, 97, 116, 101, 66, 97, 108, 97, 110, 99, 101, 34, 58, 123, 34, 100, 97, 116, 101, 34, 58, 34, 50, 48, 50, 50, 47, 49, 48, 47, 49, 53, 34, 44, 34, 97, 109, 111, 117, 110, 116, 34, 58, 49, 48, 48, 48, 125, 125, 125}

	unmarshalWrappedBody(data, "createBalance", &result)

	if result.Amount != 1000 || result.Date != "2022/10/15" {
		t.Fatalf("Failed")
	}
}
