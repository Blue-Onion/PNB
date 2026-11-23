package main

import (
	"flag"

	decoder "github.com/Blue-Onion/PNB/internal/Decoder"
	permissionchecker "github.com/Blue-Onion/PNB/internal/PermissionChecker"
)

func main() {
	hexPtr := flag.String("hexcode", "", "Hex code")
	flag.Parse()

	data, err := decoder.Decoder(*hexPtr)
	if err != nil {
		panic(err.Error())
	}
	permissionchecker.PermissionChecker(data)
}
