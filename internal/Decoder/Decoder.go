package decoder

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/Blue-Onion/PNB/internal/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// selectCall validates calldata and extracts the selector.
func selectCall(input string) (string, []byte, error) {
	if !strings.HasPrefix(input, "0x") {
		return "", nil, fmt.Errorf("missing 0x prefix")
	}

	data, err := hex.DecodeString(input[2:])
	if err != nil {
		return "", nil, fmt.Errorf("invalid hex: %w", err)
	}

	if len(data) < 4 {
		return "", nil, fmt.Errorf("calldata too short")
	}

	selector := hex.EncodeToString(data[:4])
	params := data[4:]

	return selector, params, nil
}

// decodeArgs decodes two static ABI parameters.
func decodeArgs(
	params []byte,
	firstType string,
	secondType string,
) ([]interface{}, error) {

	if len(params) != 64 {
		return nil, fmt.Errorf(
			"expected 64 parameter bytes, got %d",
			len(params),
		)
	}

	t1, err := abi.NewType(firstType, "", nil)
	if err != nil {
		return nil, err
	}

	t2, err := abi.NewType(secondType, "", nil)
	if err != nil {
		return nil, err
	}

	args := abi.Arguments{
		{Type: t1},
		{Type: t2},
	}

	return args.Unpack(params)
}

// decodeTransfer handles transfer(address,uint256).
func decodeTransfer(params []byte) (*models.Result, error) {
	values, err := decodeArgs(params, "address", "uint256")
	if err != nil {
		return nil, err
	}

	recipient := values[0].(common.Address)
	amount := values[1].(*big.Int)

	return &models.Result{
		Function: "transfer",
		Address:  recipient.Hex(),
		Amount:   amount.String(),
		Risk:     string(RiskReview),
		Warning:  WarnVerifyTransfer,
	}, nil
}

// decodeApprove handles approve(address,uint256).
func decodeApprove(params []byte) (*models.Result, error) {
	values, err := decodeArgs(params, "address", "uint256")
	if err != nil {
		return nil, err
	}

	spender := values[0].(common.Address)
	amount := values[1].(*big.Int)

	res := &models.Result{
		Function: "approve",
		Address:  spender.Hex(),
		Amount:   amount.String(),
		Risk:     string(RiskReview),
		Warning:  WarnTokenApproval,
	}

	if amount.Cmp(MaxUint256) == 0 {
		res.Risk = string(RiskHigh)
		res.Warning = WarnUnlimitedApproval
	} else if amount.Cmp(big.NewInt(maxToken)) > 0 {
		res.Risk = string(RiskReview)
		res.Warning = WarnTokenApproval
	}
	return res, nil
}

// decodeNFTApproval handles setApprovalForAll(address,bool).
func decodeNFTApproval(params []byte) (*models.Result, error) {
	values, err := decodeArgs(params, "address", "bool")
	if err != nil {
		return nil, err
	}

	operator := values[0].(common.Address)
	approved := values[1].(bool)

	res := &models.Result{
		Function: "setApprovalForAll",
		Address:  operator.Hex(),
		Risk:     string(RiskReview),
		Warning:  WarnNFTApprovalRemove,
	}

	if approved {
		res.Risk = string(RiskHigh)
		res.Warning = WarnNFTApproval
	}

	return res, nil
}

// Decoder is the main entry point.
func Decoder(input string) (*models.Result, error) {
	selector, params, err := selectCall(input)
	if err != nil {
		return nil, err
	}

	switch selector {
	case SelectorTransfer:
		return decodeTransfer(params)

	case SelectorApprove:
		return decodeApprove(params)

	case SelectorSetApprovalForAll:
		return decodeNFTApproval(params)

	default:
		return &models.Result{
			Function: "unknown",
			Risk:     string(RiskUnknown),
			Warning:  WarnUnknownFunction,
		}, nil
	}
}
