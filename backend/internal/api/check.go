package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	aichecker "github.com/Blue-Onion/txShield/internal/AiChecker"
	decoder "github.com/Blue-Onion/txShield/internal/Decoder"
	permissionchecker "github.com/Blue-Onion/txShield/internal/PermissionChecker"
	report "github.com/Blue-Onion/txShield/internal/Report"
	scamchecker "github.com/Blue-Onion/txShield/internal/ScamChecker"
	"github.com/Blue-Onion/txShield/internal/config"
	"github.com/google/uuid"
)

var errInvalidCalldata = errors.New("invalid calldata")

// checkCode runs the transaction checks and returns the JSON report.
func checkCode(logger *slog.Logger, hexCode string) ([]byte, error) {
	res, err := decoder.Decoder(hexCode)
	if err != nil {
		logErr := err
		var invalidByte hex.InvalidByteError
		if errors.As(err, &invalidByte) {
			logErr = errors.New("calldata contains a non-hexadecimal character")
		}
		logger.Warn("Decode calldata", "error", logErr)
		return nil, fmt.Errorf("%w: %w", errInvalidCalldata, err)
	}

	permissionchecker.PermissionChecker(res)
	if err := scamchecker.ScamAddressChecker(res); err != nil {
		logger.Error("Check scam addresses", "error", err)
	}
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Load configuration", "error", err)
		return []byte(""), err
	}
	ai, err := aichecker.Check(res, cfg.GeminiAPIKey, logger)

	reportID := uuid.New().String()
	result := report.Build(reportID, res, ai)
	if err != nil {
		result.Analysis.Status = "AI_UNAVAILABLE"
	}
	data, marshalErr := json.MarshalIndent(result, "", "  ")
	if marshalErr != nil {
		logger.Error("Encode transaction report", "error", marshalErr)
	}
	return data, marshalErr
}
