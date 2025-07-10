package addressutils

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"github.com/pariz/gountries"
)

func CountryCodeToName(code string) string {
	region, err := language.ParseRegion(code)
	if err != nil {
		return ""
	}

	return display.English.Regions().Name(region)
}

func IsValidCountryCode(code string) bool {
	parsedCode, err := language.ParseRegion(code)
	if err != nil {
		return false
	}

	return parsedCode.IsCountry()
}

func CountryNameToCode(name string) string {
	query := gountries.New()
	country, err := query.FindCountryByName(name)
	if err != nil {
		return ""
	}
	return country.Alpha2
}

func LinkedInLocationToCityState(location string) (city string, state string) {
	parts := strings.Split(location, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	switch len(parts) {
	case 3:
		city = parts[0]
		state = parts[1]
	case 2:
		state = parts[0]
	}

	return city, state
}
