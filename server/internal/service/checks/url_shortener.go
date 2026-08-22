package checks

import "github.com/AARYA-001/the-url-mentalist/internal/constants"

func IsUrlShortener(domain string) bool {
	_, ok := constants.URLShorteners[domain]
	return ok
}
