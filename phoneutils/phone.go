package phoneutils

import (
	"regexp"
	"strings"
)

func FormatPhonesAndProviders(providers, phones []string) string {
	var result []string

	for i, phone := range phones {
		if i < len(providers) {
			result = append(result, providers[i]+": "+phone)
		} else {
			result = append(result, phone)
		}
	}

	return strings.Join(result, ", ")
}

// Clean formats a phone number string by:
// - Preserving a leading + if present
// - Removing all other + symbols
// - Removing all non-digit characters (spaces, dashes, parentheses, etc)
// Examples:
//
//	"+1 (555) 123-4567" -> "+15551234567"
//	"1+555-123-4567"    -> "15551234567"
func Clean(phone string) string {
	cleaned := regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")
	if cleaned == "" {
		return ""
	}

	if strings.HasPrefix(phone, "+") {
		return "+" + cleaned
	}

	return cleaned
}
