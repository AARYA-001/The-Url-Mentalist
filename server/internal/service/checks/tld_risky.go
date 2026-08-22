package checks

import "github.com/AARYA-001/the-url-mentalist/internal/constants"

func IsRiskyTld(domain string) (bool, bool, string) {
	tld, icann := GetTld(domain)
	_, ok := constants.RiskyTLDs[tld]

	return ok, icann, tld
}
