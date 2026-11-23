package aichecker

type AiReport struct {
	RiskPoints      []RiskPoint `json:"risk_points"`
	Summary         string      `json:"summary"`
	Recommendations []string    `json:"recommendations"`
	Confidence      string      `json:"confidence"`
} // AI output structure
type RiskPoint struct {
	Severity    string `json:"severity"`
	Reason      string `json:"reason"`
	Explanation string `json:"explanation"`
}
