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

func TestGetCountryAlpha2FromLocation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		location      string
		expectedFound bool
		expected      string
	}

	tests := []testCase{
		{
			name:          "should return empty string and false when location is empty",
			location:      "",
			expectedFound: false,
			expected:      "",
		},
		{
			name:          "should return empty string and false when location does not contain a country",
			location:      "Some random text",
			expectedFound: false,
			expected:      "",
		},
		{
			name:          "should return the country alpha-2 code string and true when location contains a country",
			location:      "Cairo, Egypt",
			expectedFound: true,
			expected:      "EG",
		},
		{
			name:          "should return the country alpha-2 code string and true when location contains a country separated by commas and no spaces",
			location:      "Chicago,Illinois,United States",
			expectedFound: true,
			expected:      "US",
		},
		{
			name:          "should return the country alpha-2 code string and true when location contains a country separated by spaces and no commas",
			location:      "Chicago Illinois United States",
			expectedFound: true,
			expected:      "US",
		},
		{
			name:          "should return the country alpha-2 code of the latest country name in the string and true when location contains two places that are country names",
			location:      "Cairo, New York, United States",
			expectedFound: true,
			expected:      "US",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := GetCountryAlpha2FromLocation(tt.location)
			require.Equal(t, tt.expectedFound, found)
			require.Equal(t, tt.expected, got)
		})
	}
}
