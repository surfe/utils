package addressutils

import (
	"testing"
)

func TestCodeToCountryName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "Valid country code - France",
			code:     "FR",
			expected: "France",
		},
		{
			name:     "Valid country code - United States",
			code:     "US",
			expected: "United States",
		},
		{
			name:     "Valid country code - Japan",
			code:     "JP",
			expected: "Japan",
		},
		{
			name:     "Valid country code - Germany",
			code:     "DE",
			expected: "Germany",
		},
		{
			name:     "Empty string",
			code:     "",
			expected: "",
		},
		{
			name:     "Invalid country code",
			code:     "XX",
			expected: "",
		},
		{
			name:     "Lowercase code",
			code:     "fr",
			expected: "France",
		},
		{
			name:     "Numbers instead of letters",
			code:     "12",
			expected: "",
		},
		{
			name:     "Too long code",
			code:     "FRA",
			expected: "France",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := CountryCodeToName(tt.code)
			if result != tt.expected {
				t.Errorf("CountryCodeToName(%q) = %q, want %q",
					tt.code, result, tt.expected)
			}
		})
	}
}

func TestLinkedInLocationToCityState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		location  string
		wantCity  string
		wantState string
	}{
		// Three part locations (City, State, Country)
		{
			name:      "San Jose CA",
			location:  "San Jose, California, United States",
			wantCity:  "San Jose",
			wantState: "California",
		},
		{
			name:      "North York Ontario",
			location:  "North York, Ontario, Canada",
			wantCity:  "North York",
			wantState: "Ontario",
		},
		{
			name:      "Ankara Turkey",
			location:  "Ankara, Ankara, Türkiye",
			wantCity:  "Ankara",
			wantState: "Ankara",
		},
		{
			name:      "Cuauhtémoc Mexico",
			location:  "Cuauhtémoc, Chihuahua, Mexico",
			wantCity:  "Cuauhtémoc",
			wantState: "Chihuahua",
		},
		{
			name:      "Vienna Austria",
			location:  "Vienna, Vienna, Austria",
			wantCity:  "Vienna",
			wantState: "Vienna",
		},

		// Two part locations (Region/City, Country)
		{
			name:      "Delhi India",
			location:  "Delhi, India",
			wantCity:  "",
			wantState: "Delhi",
		},
		{
			name:      "Budapest Hungary",
			location:  "Budapest, Hungary",
			wantCity:  "",
			wantState: "Budapest",
		},
		{
			name:      "Paris France",
			location:  "Paris, France",
			wantCity:  "",
			wantState: "Paris",
		},

		// Metropolitan areas and special formats
		{
			name:      "Greater Lyon Area",
			location:  "Greater Lyon Area",
			wantCity:  "",
			wantState: "",
		},
		{
			name:      "Frankfurt Metro",
			location:  "Frankfurt Rhine-Main Metropolitan Area",
			wantCity:  "",
			wantState: "",
		},

		// Edge cases
		{
			name:      "Empty string",
			location:  "",
			wantCity:  "",
			wantState: "",
		},
		{
			name:      "Single city",
			location:  "Tokyo",
			wantCity:  "",
			wantState: "",
		},
		{
			name:      "Extra spaces",
			location:  "  San Jose  ,  California  ,  United States  ",
			wantCity:  "San Jose",
			wantState: "California",
		},
		{
			name:      "Union County NJ",
			location:  "Union County, New Jersey, United States",
			wantCity:  "Union County",
			wantState: "New Jersey",
		},
		{
			name:      "Gdańsk Poland",
			location:  "Gdańsk, Pomorskie, Poland",
			wantCity:  "Gdańsk",
			wantState: "Pomorskie",
		},
		{
			name:      "Mumbai India",
			location:  "Mumbai, Maharashtra, India",
			wantCity:  "Mumbai",
			wantState: "Maharashtra",
		},
		{
			name:      "Ville de Paris",
			location:  "Ville de Paris, Île-de-France, France",
			wantCity:  "Ville de Paris",
			wantState: "Île-de-France",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotCity, gotState := LinkedInLocationToCityState(tt.location)
			if gotCity != tt.wantCity {
				t.Errorf("LinkedInLocationToCityState() city = %v, want %v", gotCity, tt.wantCity)
			}

			if gotState != tt.wantState {
				t.Errorf("LinkedInLocationToCityState() state = %v, want %v", gotState, tt.wantState)
			}
		})
	}
}
