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

const (
	SelectorTransfer          = "a9059cbb"
	SelectorApprove           = "095ea7b3"
	SelectorSetApprovalForAll = "a22cb465"
)

func selectCall(input string) (selector string, params []byte, err error) {
	if !strings.HasPrefix(input, "0x") {
		return "", nil, fmt.Errorf("Malformed Hexcode")
	}

	data, err := hex.DecodeString(input[2:])
	if err != nil {
		return "", nil, err
	}
	if len(data) < 4 {
		return "", nil, fmt.Errorf("Call Data too Short")
	}
	return hex.EncodeToString(data[:4]), data[4:], nil
}

func Decoder(input string) (*models.Result, error) {
	selector, params, err := selectCall(input)
	if err != nil {
		return nil, err
	}
	res := models.Result{}
	addressType, _ := abi.NewType("address", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	boolType, _ := abi.NewType("bool", "", nil)
	switch selector {
	case SelectorApprove:
		args := abi.Arguments{
			{Type: addressType},
			{Type: uintType},
		}
		values, err := args.Unpack(params)
		if err != nil {
			return &res, err
		}
		spendor := values[0].(common.Address)
		amount := values[1].(*big.Int)
		maxUnit := new(big.Int).Sub(
			new(big.Int).Lsh(big.NewInt(1), 256),
			big.NewInt(1),
		)
		res.Function = "approve"
		res.Address = spendor.Hex()
		res.Amount = amount.String()
		if amount.Cmp(maxUnit) == 0 {
			res.Risk = "HIGH"
			res.Warning = "Unlimited Token Spending approval"
		} else {
			res.Risk = "REVIEW"
			res.Warning = "Token spending permission requested"
		}
	case SelectorTransfer:
		args := abi.Arguments{
			{Type: addressType},
			{Type: uintType},
		}
		values, err := args.Unpack(params)
		if err != nil {
			return &res, err
		}
		recipent := values[0].(common.Address)
		amount := values[1].(*big.Int)
		res.Function = "transfer"
		res.Address = recipent.Hex()
		res.Amount = amount.String()
		res.Risk = "REVIEW"
		res.Warning = "Verify It first"
	case SelectorSetApprovalForAll:

		args := abi.Arguments{
			{Type: addressType},
			{Type: boolType},
		}
		values, err := args.Unpack(params)
		if err != nil {
			return &res, err
		}
		opr := values[0].(common.Address)
		apr := values[1].(bool)
		res.Function = "setApprovalForAll"
		res.Address = opr.Hex()

		if apr {
			res.Risk = "HIGH"
			res.Warning = "Collection-wide NFT operator approval"
		} else {
			res.Risk = "REVIEW"
			res.Warning = "NFT operator approval removal requested"
		}
	default:
		res.Function = "unknown"
		res.Risk = "UNKNOWN"
		res.Warning = "Unrecognized function selector"
	}
	return &res, nil
}
