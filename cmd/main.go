package main

import (
	"flag"
	"fmt"

	decoder "github.com/Blue-Onion/PNB/internal/Decoder"
)

func check(s *string) {
	fmt.Println(*s)
}
func main() {
	hexPtr := flag.String("hexcode", "", "Hex code")
	flag.Parse()

	fmt.Println(*hexPtr)
	check(hexPtr)
	data, err := decoder.Decoder(*hexPtr)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(data)
}
