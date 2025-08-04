package utils

import (
	"context"
	"fmt"
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

// GetCountryAlpha2FromLocation finds the country within a location string and returns:
// 1. The Alpha 2 ISO code for that country
// 2. a boolean that denotes whether or not the country was found
// Examples:
// - input: "Cairo, Egypt" > output: "EG", true
// - input: "Some random text" > output: "", false
func GetCountryAlpha2FromLocation(location string) (string, bool) {
	var combinedParts string
	location = strings.ReplaceAll(location, ",", " ")
	locationParts := strings.Split(location, " ")
	for _, v := range slices.Backward(locationParts) {
		v = strings.TrimSpace(v)
		alpha2 := GetCountryAlpha2(context.Background(), "", v)
		if alpha2 != "" {
			return alpha2, true
		}

		combinedParts = strings.TrimSpace(fmt.Sprintf("%s %s", v, combinedParts))
		alpha2 = GetCountryAlpha2(context.Background(), "", combinedParts)
		if alpha2 != "" {
			return alpha2, true
		}
	}

	return "", false
}
