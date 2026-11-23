package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	decoder "github.com/Blue-Onion/PNB/internal/Decoder"
	permissionchecker "github.com/Blue-Onion/PNB/internal/PermissionChecker"
	scamchecker "github.com/Blue-Onion/PNB/internal/ScamChecker"
)

func main() {
	hexPtr := flag.String("hexcode", "", "Hex code")
	flag.Parse()

	res, err := decoder.Decoder(*hexPtr)
	if err != nil {
		panic(err.Error())
	}
	permissionchecker.PermissionChecker(res)
	if err := scamchecker.ScamAddressChecker(res); err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(res); err != nil {
		log.Fatal(err)
	}
}
