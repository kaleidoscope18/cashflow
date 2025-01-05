package bdd

import (
	"cashflow/models"
	"testing"
)

func TestUnmarshalWrappedBody(t *testing.T) {
	var result models.Balance

	data := []byte(`{"data":{"createBalance":{"date":"2022-10-15","amount":1000}}}`)

	unmarshalWrappedBody(data, "createBalance", &result)

	if result.Amount != 1000 || result.Date != "2022-10-15" {
		t.Fatalf("Failed: %s", result.Date)
	}
}
