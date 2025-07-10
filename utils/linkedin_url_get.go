package utils

type linkedinURLGetter interface {
	GetLinkedinURL(key string) string
}

func RemoveWithEmptyOrNotEqualLinkedinURL[T linkedinURLGetter](entities []T, key, linkedinURL string) []T {
	cleanLinkedInURL := LinkedinURLCleaner(linkedinURL)
	if cleanLinkedInURL == "" {
		return entities
	}

	remainingRecords := make([]T, 0)
	for _, e := range entities {
		crmLinkedinURL := e.GetLinkedinURL(key)
		cleanCRMLinkedinURL := LinkedinURLCleaner(crmLinkedinURL)
		if cleanCRMLinkedinURL == "" || ExtractHostAndPath(cleanCRMLinkedinURL) == ExtractHostAndPath(cleanLinkedInURL) {
			remainingRecords = append(remainingRecords, e)
		}
	}

	return remainingRecords
}
