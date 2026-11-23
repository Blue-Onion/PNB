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

// Decoder validates calldata and extracts the arguments for recognized functions.
func Decoder(input string) (*models.Result, error) {
	selector, params, err := selectCall(input)
	if err != nil {
		return nil, err
	}

	res := &models.Result{Findings: []models.Finding{}}
	switch selector {
	case SelectorTransfer, SelectorApprove:
		values, err := decodeArgs(params, "address", "uint256")
		if err != nil {
			return nil, err
		}

		res.Function = "transfer"
		if selector == SelectorApprove {
			res.Function = "approve"
		}
		res.Address = values[0].(common.Address).Hex()
		res.Amount = values[1].(*big.Int).String()

	case SelectorSetApprovalForAll:
		values, err := decodeArgs(params, "address", "bool")
		if err != nil {
			return nil, err
		}

		approved := values[1].(bool)
		res.Function = "setApprovalForAll"
		res.Address = values[0].(common.Address).Hex()
		res.Approved = &approved

	default:
		res.Function = "unknown"
	}
	return res, nil
}
