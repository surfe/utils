package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCountryAlphaCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		code     string
		expected bool
	}{
		{"US", true},
		{"GB", true},
		{"xyz", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			t.Parallel()

			result := IsCountryAlphaCode(tt.code)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCountryNameToCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected string
	}{
		{"United States", "US"},
		{"united states", "US"},
		{"Germany", "DE"},
		{"Nonexistentland", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := CountryNameToCode(tt.name)
			require.Equal(t, tt.expected, got)
		})
	}
}
