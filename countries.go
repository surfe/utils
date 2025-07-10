package utils

import (
	"slices"
	"strings"

	"github.com/pariz/gountries"
)

var euCountryCodes = []string{"at", "be", "bg", "hr", "cy", "cz", "dk", "ee", "fi", "fr", "de", "gr", "hu", "ie", "it", "lv", "lt", "lu", "mt", "nl", "pl", "pt", "ro", "sk", "si", "es", "se"}

func CountryNameToCode(name string) string {
	query := gountries.New()

	country, err := query.FindCountryByName(name)
	if err != nil {
		return ""
	}

	return country.Alpha2
}

func IsCountryAlphaCode(country string) bool {
	query := gountries.New()
	_, err := query.FindCountryByAlpha(country)

	return err == nil
}

func IsCountryInEU(country string) bool {
	if !IsCountryAlphaCode(country) {
		country = CountryNameToCode(country)
		if country == "" {
			return false
		}
	}

	return slices.Contains(euCountryCodes, strings.ToLower(country))
}
