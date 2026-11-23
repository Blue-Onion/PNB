package permissionchecker

import (
	"math/big"

	"github.com/Blue-Onion/txShield/internal/models"
)

// Risk levels
type RiskLevel string

const (
	RiskHigh    RiskLevel = "HIGH"
	RiskReview  RiskLevel = "REVIEW"
	RiskUnknown RiskLevel = "UNKNOWN"
)

// Reusable warning messages
const (
	WarnUnlimitedApproval  = "Unlimited token spending approval"
	WarnTokenApproval      = "Token spending permission requested"
	WarnLargeTokenApproval = "Token approval exceeds 50,000 raw units. If this is a free mint, verify what the token approval is for before signing."
	WarnVerifyTransfer     = "Verify recipient and token amount"
	WarnNFTApproval        = "Collection-wide NFT operator approval"
	WarnNFTApprovalRemove  = "NFT operator approval removal requested"
	WarnUnknownFunction    = "Unrecognized function selector"
)

const (
	tokenApprovalThreshold = 50_000
	engineName             = "PERMISSION_CHECKER"
)

var maxUint256 = new(big.Int).Sub(
	new(big.Int).Lsh(big.NewInt(1), 256),
	big.NewInt(1),
)

// PermissionChecker adds the permission assessment to res.
func PermissionChecker(res *models.Result) {
	var severity models.RiskLevel
	var message string
	switch res.Function {
	case "transfer":
		res.Risk = string(RiskReview)
		res.Warning = WarnVerifyTransfer
	case "approve":
		res.Risk = string(RiskReview)
		res.Warning = WarnTokenApproval
		amount, ok := new(big.Int).SetString(res.Amount, 10)
		if ok && amount.Cmp(maxUint256) == 0 {
			res.Risk = string(RiskHigh)
			res.Warning = WarnUnlimitedApproval
		} else if ok && amount.Cmp(big.NewInt(tokenApprovalThreshold)) > 0 {
			res.Warning = WarnLargeTokenApproval
			severity = models.RiskMedium
			message = "Approval exceeds 50000 raw token units"
		}
	case "setApprovalForAll":
		res.Risk = string(RiskReview)
		res.Warning = WarnNFTApprovalRemove
		if res.Approved != nil && *res.Approved {
			res.Risk = string(RiskHigh)
			res.Warning = WarnNFTApproval
		}
	default:
		res.Risk = string(RiskUnknown)
		res.Warning = WarnUnknownFunction
	}
	if severity == "" {
		severity = models.RiskLevel(res.Risk)
		message = res.Warning
	}
	res.Findings = append(res.Findings, models.Finding{
		Engine: engineName, Severity: severity, Message: message,
	})
}
