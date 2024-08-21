package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type (
	Trade struct {
		Symbol       string
		FirstAmount  float64 `json:"first_amount"`
		SecondAmount float64 `json:"second_amount"`
	}
)

func ParseTrades(file string) []Trade {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open trades file")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer jsonFile.Close()

	var trades []Trade

	byteValue, _ := io.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &trades)

	return trades
}
