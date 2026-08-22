package constants

import "os"

const DOMAIN_RANK_FILE_PATH = "./assets/top-1m.csv"

var DomainRankCandidatePaths = []string{
	"./assets/top-1m.csv",
	"server/assets/top-1m.csv",
	"../assets/top-1m.csv",
	"/app/assets/top-1m.csv",
}

// FindDomainRankFilePath resolves the top-1m.csv location across different working directories.
func FindDomainRankFilePath() string {
	for _, p := range DomainRankCandidatePaths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return DOMAIN_RANK_FILE_PATH
}
