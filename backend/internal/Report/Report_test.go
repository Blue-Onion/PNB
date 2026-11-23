package report

import (
	"encoding/json"
	"strings"
	"testing"

	aichecker "github.com/Blue-Onion/txShield/internal/AiChecker"
	decoder "github.com/Blue-Onion/txShield/internal/Decoder"
	permissionchecker "github.com/Blue-Onion/txShield/internal/PermissionChecker"
	scamchecker "github.com/Blue-Onion/txShield/internal/ScamChecker"
	"github.com/Blue-Onion/txShield/internal/models"
)

func TestBuildReportForFlaggedApproval(t *testing.T) {
	calldata := "0x095ea7b3" + strings.Repeat("0", 24) + strings.Repeat("1", 40) + strings.Repeat("0", 60) + "c351"
	res, err := decoder.Decoder(calldata)
	if err != nil {
		t.Fatal(err)
	}
	permissionchecker.PermissionChecker(res)
	if err := scamchecker.ScamAddressChecker(res); err != nil {
		t.Fatal(err)
	}

	ai := aichecker.AiReport{Summary: "Review the requested allowance.", Confidence: "MEDIUM"}
	report := Build("txShield-001", res, ai)
	if report.FinalRisk != models.RiskHigh || report.Transaction.Amount != "50001" {
		t.Fatalf("unexpected risk or transaction: %+v", report)
	}
	if len(report.Findings) != 2 || report.Findings[0].Severity != models.RiskMedium || report.Findings[1].Severity != models.RiskHigh {
		t.Fatalf("unexpected findings: %+v", report.Findings)
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"report_id", "transaction", "final_risk", "findings", "ai_analysis", "analysis"} {
		if _, ok := fields[field]; !ok {
			t.Errorf("report JSON is missing %q", field)
		}
	}
}
