package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	aichecker "github.com/Blue-Onion/PNB/internal/AiChecker"
	decoder "github.com/Blue-Onion/PNB/internal/Decoder"
	permissionchecker "github.com/Blue-Onion/PNB/internal/PermissionChecker"
	report "github.com/Blue-Onion/PNB/internal/Report"
	scamchecker "github.com/Blue-Onion/PNB/internal/ScamChecker"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	hexPtr := flag.String("hexcode", "", "Hex code")
	reportID := flag.String("report-id", "PNB-001", "Report identifier")
	flag.Parse()

	res, err := decoder.Decoder(*hexPtr)
	if err != nil {
		panic(err.Error())
	}
	permissionchecker.PermissionChecker(res)
	if err := scamchecker.ScamAddressChecker(res); err != nil {
		log.Println(err)
	}

	ai, err := aichecker.Check(res)
	if err != nil {
		log.Println(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report.Build(*reportID, res, ai)); err != nil {
		log.Println(err)
	}
}
