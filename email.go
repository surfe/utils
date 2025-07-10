package utils

import "net/mail"

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsPersonalEmail(email string) bool {
	if email == "" {
		return false
	}

	domain := DomainFromEmail(email)
	_, found := personalEmailProvidersDomains[domain]

	return found
}

func IsPersonalEmailDomain(domain string) bool {
	_, found := personalEmailProvidersDomains[domain]

	return found
}
