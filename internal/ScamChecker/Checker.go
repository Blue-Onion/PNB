package scamchecker

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/Blue-Onion/PNB/internal/models"
)

const (
	EngineName    = "SCAM_ADDRESS_CHECKER"
	WarnKnownScam = "Address found in local scam database"
)

func ScamAddressChecker(
	res *models.Result,
) error {
	if res == nil {
		return fmt.Errorf("Empty Res")
	}
	if res.Address == "" {
		return nil
	}
	isValid := checkDb(res.Address)
	if isValid {
		res.Findings = append(res.Findings, models.Finding{
			Engine:   EngineName,
			Severity: models.RiskHigh,
			Message:  WarnKnownScam,
		})
	}
	return nil
}

//go:embed scam.txt
var scamAddresses string

func checkDb(address string) bool {
	for _, entry := range strings.Split(scamAddresses, "\n") {
		if strings.EqualFold(strings.TrimSpace(entry), address) {
			return true
		}
	}
	return false
}
