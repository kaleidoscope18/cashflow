package graph

import (
	"testing"
)

func TestValidateRecurrency(t *testing.T) {
	validRecurrency := "FREQ=DAILY;INTERVAL=1"
	invalidRecurrency := "INVALID"

	tests := []struct {
		name        string
		recurrency  *string
		expectError bool
	}{
		{
			name:        "Valid recurrency",
			recurrency:  &validRecurrency,
			expectError: false,
		},
		{
			name:        "Invalid recurrency",
			recurrency:  &invalidRecurrency,
			expectError: true,
		},
		{
			name:        "Nil recurrency",
			recurrency:  nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecurrency(tt.recurrency)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
