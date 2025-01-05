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
		expected    string
		expectError bool
	}{
		{
			name:        "Valid recurrency",
			recurrency:  &validRecurrency,
			expected:    "FREQ=DAILY;INTERVAL=1",
			expectError: false,
		},
		{
			name:        "Invalid recurrency",
			recurrency:  &invalidRecurrency,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Nil recurrency",
			recurrency:  nil,
			expected:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateRecurrency(tt.recurrency)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
