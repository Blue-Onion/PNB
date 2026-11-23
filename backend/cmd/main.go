package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	aichecker "github.com/Blue-Onion/PNB/internal/AiChecker"
	decoder "github.com/Blue-Onion/PNB/internal/Decoder"
	permissionchecker "github.com/Blue-Onion/PNB/internal/PermissionChecker"
	report "github.com/Blue-Onion/PNB/internal/Report"
	scamchecker "github.com/Blue-Onion/PNB/internal/ScamChecker"
	"github.com/Blue-Onion/PNB/internal/config"
)

func main() {
	hexPtr := flag.String("hexcode", "", "Hex code")
	reportID := flag.String("report-id", "PNB-001", "Report identifier")
	serve := flag.Bool("serve", false, "Start the local API server")
	addr := flag.String("addr", "", "Override the API listen address (default localhost:PORT)")
	cert := flag.String("tls-cert", "", "TLS certificate file")
	key := flag.String("tls-key", "", "TLS private key file")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	check := func(hexCode, reportID string) ([]byte, error) {
		return checkCode(hexCode, reportID, cfg.GeminiAPIKey)
	}

	if *serve {
		listenAddr := *addr
		if listenAddr == "" {
			listenAddr = net.JoinHostPort("localhost", cfg.Port)
		}
		log.Fatal(serveAPI(listenAddr, *cert, *key, *reportID, check))
	}

	data, err := check(*hexPtr, *reportID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}

var errInvalidCalldata = errors.New("invalid calldata")

// checkCode runs the transaction checks and returns the JSON report.
func checkCode(hexCode, reportID, geminiAPIKey string) ([]byte, error) {
	res, err := decoder.Decoder(hexCode)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errInvalidCalldata, err)
	}

	permissionchecker.PermissionChecker(res)
	if err := scamchecker.ScamAddressChecker(res); err != nil {
		log.Println(err)
	}

	ai, err := aichecker.Check(res, geminiAPIKey)
	if err != nil {
		log.Println(err)
	}

	result := report.Build(reportID, res, ai)
	if err != nil {
		result.Analysis.Status = "AI_UNAVAILABLE"
	}
	return json.MarshalIndent(result, "", "  ")
}
