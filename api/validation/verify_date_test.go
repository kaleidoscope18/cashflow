package validation

import (
	"testing"
)

func TestVerifyDateHappyPath(t *testing.T) {
	validDate := "2025-01-01"
	result, err := VerifyDate(&validDate)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if *result != validDate {
		t.Errorf("Expected %v, got %v", validDate, result)
	}
}

func TestVerifyDateOtherCases(t *testing.T) {
	invalidDate := "INVALID"
	unexistentDate := "2025-01-32"
	var nilDate *string
	emptyDate := ""

	tests := []struct {
		name        string
		date        *string
		expected    *string
		expectError bool
	}{
		{
			name:        "Invalid date",
			date:        &invalidDate,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Unexistent date",
			date:        &unexistentDate,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Nil date",
			date:        nilDate,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "Empty date",
			date:        &emptyDate,
			expected:    nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := VerifyDate(tt.date)
			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
