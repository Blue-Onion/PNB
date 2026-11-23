package report

import (
	aichecker "github.com/Blue-Onion/PNB/internal/AiChecker"
	"github.com/Blue-Onion/PNB/internal/models"
)

type Transaction struct {
	Function string `json:"function"`
	Address  string `json:"address,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Approved *bool  `json:"approved,omitempty"`
}

type Analysis struct {
	SimulationPerformed bool   `json:"simulation_performed"`
	ReputationSource    string `json:"reputation_source"`
	Status              string `json:"status"`
}

type Report struct {
	ID          string             `json:"report_id"`
	Transaction Transaction        `json:"transaction"`
	FinalRisk   models.RiskLevel   `json:"final_risk"`
	Findings    []models.Finding   `json:"findings"`
	AIAnalysis  aichecker.AiReport `json:"ai_analysis"`
	Analysis    Analysis           `json:"analysis"`
}

// Build combines a completed local assessment with the AI analysis.
func Build(id string, res *models.Result, ai aichecker.AiReport) Report {
	riskOrder := map[models.RiskLevel]int{
		models.RiskUnknown: 0,
		models.RiskSafe:    1,
		models.RiskLow:     2,
		models.RiskReview:  3,
		models.RiskMedium:  4,
		models.RiskHigh:    5,
	}
	risk := models.RiskLevel(res.Risk)
	if _, ok := riskOrder[risk]; !ok {
		risk = models.RiskUnknown
	}
	for _, finding := range res.Findings {
		if riskOrder[finding.Severity] > riskOrder[risk] {
			risk = finding.Severity
		}
	}

	return Report{
		ID: id,
		Transaction: Transaction{
			Function: res.Function,
			Address:  res.Address,
			Amount:   res.Amount,
			Approved: res.Approved,
		},
		FinalRisk:  risk,
		Findings:   append([]models.Finding{}, res.Findings...),
		AIAnalysis: ai,
		Analysis: Analysis{
			ReputationSource: "LOCAL_TEST_DATABASE",
			Status:           "COMPLETED",
		},
	}
}
