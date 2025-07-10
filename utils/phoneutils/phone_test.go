package phoneutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPhonesAndProviders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		providers []string
		phones    []string
		wanted    string
	}{
		{
			name:      "Given single provider and single phone should return expected string",
			providers: []string{"Provider1"},
			phones:    []string{"1234567890"},
			wanted:    "Provider1: 1234567890",
		},
		{
			name:      "Given multiple providers and multiple phones should return expected string",
			providers: []string{"Provider1", "Provider2"},
			phones:    []string{"1234567890", "0987654321"},
			wanted:    "Provider1: 1234567890, Provider2: 0987654321",
		},
		{
			name:      "Given multiple providers and single phone should return expected string",
			providers: []string{"Provider1"},
			phones:    []string{"1234567890", "0987654321"},
			wanted:    "Provider1: 1234567890, 0987654321",
		},
		{
			name:      "Given single provider and multiple phones should return expected string",
			providers: []string{"Provider1", "Provider2"},
			phones:    []string{"1234567890"},
			wanted:    "Provider1: 1234567890",
		},
		{
			name:   "Given no providers and single phone should return expected string",
			phones: []string{"1234567890"},
			wanted: "1234567890",
		},
		{
			name:      "Given no providers and multiple phones should return expected string",
			providers: []string{"Provider1"},
			wanted:    "",
		},
		{
			name:   "Given no providers and no phones should return expected string",
			wanted: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := FormatPhonesAndProviders(tt.providers, tt.phones)
			assert.Equal(t, tt.wanted, actual)
		})
	}
}

func TestClean(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "international format with leading plus",
			input:    "+1 (555) 123-4567",
			expected: "+15551234567",
		},
		{
			name:     "local format without plus",
			input:    "555-123-4567",
			expected: "5551234567",
		},
		{
			name:     "plus in middle",
			input:    "1+555-123-4567",
			expected: "15551234567",
		},
		{
			name:     "multiple plus signs with leading plus",
			input:    "++1-555-123+4567",
			expected: "+15551234567",
		},
		{
			name:     "multiple plus signs without leading plus",
			input:    "1++555+123+4567",
			expected: "15551234567",
		},
		{
			name:     "various separators",
			input:    "+1 (555) 123.4567",
			expected: "+15551234567",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only non-digits",
			input:    "+-() .",
			expected: "",
		},
		{
			name:     "only plus signs",
			input:    "+++",
			expected: "",
		},
		{
			name:     "text only",
			input:    "not found",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Clean(tt.input)
			if got != tt.expected {
				t.Errorf("FormatPhoneNumber(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
