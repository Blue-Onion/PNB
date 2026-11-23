package decoder

import "math/big"

// Function selectors
const (
	SelectorTransfer          = "a9059cbb"
	SelectorApprove           = "095ea7b3"
	SelectorSetApprovalForAll = "a22cb465"
)

// Risk levels
type RiskLevel string

const (
	RiskHigh    RiskLevel = "HIGH"
	RiskReview  RiskLevel = "REVIEW"
	RiskUnknown RiskLevel = "UNKNOWN"
	RiskSafe    RiskLevel = "SAFE"
)

// Reusable warning messages
const (
	WarnUnlimitedApproval = "Unlimited token spending approval"
	WarnTokenApproval     = "Token spending permission requested"
	WarnVerifyTransfer    = "Verify recipient and token amount"
	WarnNFTApproval       = "Collection-wide NFT operator approval"
	WarnNFTApprovalRemove = "NFT operator approval removal requested"
	WarnUnknownFunction   = "Unrecognized function selector"
)

var MaxUint256 = new(big.Int).Sub(
	new(big.Int).Lsh(big.NewInt(1), 256),
	big.NewInt(1),
)

const maxToken = 500000
