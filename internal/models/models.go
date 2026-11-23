package models

type Result struct {
	Function string `json:"function"`
	Address  string `json:"address,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Risk     string `json:"risk"`
	Warning  string `json:"warning"`
}
