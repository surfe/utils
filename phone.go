package utils

import "regexp"

var personalPhone = regexp.MustCompile(`^((\+|00)33|0)\s*(6|7)([\s.-]*\d{2}){4}$`)

func IsPersonalPhone(phone string) bool {
	return personalPhone.MatchString(phone)
}
