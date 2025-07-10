package utils

import "reflect"

func CountNonEmptyFields(s interface{}) int {
	val := reflect.ValueOf(s)
	count := 0

	// Ensure we have a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return 0
	}

	// Iterate over fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Skip slices with len(0)
		if field.Kind() == reflect.Slice && field.Len() == 0 {
			continue
		}

		// Check for non-zero values
		if !field.IsZero() {
			count++
		}
	}

	return count
}
