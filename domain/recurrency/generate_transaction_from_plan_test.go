package recurrency

import (
	"testing"
)

func TestGenerateTransactionFromPlan(t *testing.T) {
	transaction, err := GenerateTransactionFromPlan("ping-pong table", "chequing", "2025-01-01", "FREQ=MONTHLY;INTERVAL=2;COUNT=12", 1000)
	if err != nil {
		t.Errorf("Error generating transaction: %v", err)
	}

	if transaction.Amount != 83.33 {
		t.Fatalf("amount should be 83.33, got: %+v", transaction)
	}

	if transaction.Date != "2025-01-01" {
		t.Fatalf("date should be 2025-01-01, got: %+v", transaction)
	}

	if transaction.Recurrency != "FREQ=MONTHLY;INTERVAL=2;COUNT=12" {
		t.Fatalf("recurrency should be FREQ=MONTHLY;INTERVAL=2;COUNT=12, got: %+v", transaction)
	}

	if transaction.Account != "chequing" {
		t.Fatalf("account should be chequing, got: %+v", transaction)
	}

	if transaction.Description != "from plan: ping-pong table" {
		t.Fatalf("description should be \"from plan: ping-pong table\", got: %+v", transaction)
	}
}
