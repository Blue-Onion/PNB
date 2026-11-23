package models

type RiskLevel string

const (
	RiskHigh    RiskLevel = "HIGH"
	RiskMedium  RiskLevel = "MEDIUM"
	RiskLow     RiskLevel = "LOW"
	RiskReview  RiskLevel = "REVIEW"
	RiskUnknown RiskLevel = "UNKNOWN"
	RiskSafe    RiskLevel = "SAFE"
)

type Finding struct {
	Engine   string    `json:"engine"`
	Severity RiskLevel `json:"severity"`
	Message  string    `json:"message"`
}
type Result struct {
	Function string `json:"function"`
	Address  string `json:"address,omitempty"`
	Amount   string `json:"amount,omitempty"`

	// Existing fields (temporary compatibility)
	Risk    string `json:"risk,omitempty"`
	Warning string `json:"warning,omitempty"`

	// NFT approval status
	Approved *bool `json:"approved,omitempty"`

	// Findings from all security engines
	Findings []Finding `json:"findings"`
}
