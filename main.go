package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func help() {
	fmt.Println("Provide a configuration TOML file as the argument in CLI")
}

func dump(data interface{}) string {
	b, _ := json.MarshalIndent(data, "", "  ")
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		help()
		os.Exit(0)
	}

	arg := os.Args[1]

	if arg == "help" || arg == "-help" || arg == "--help" {
		help()
		os.Exit(0)
	}

	config := ParseConfig(arg)
	trades := ParseTrades(config.TradesFile)

	prices := GetData(config.Binance.Apikey, config.Binance.Secretkey, trades)

	fmt.Println("Latest prices:")
	fmt.Print(dump(prices))

	SendNotification(config.Telegram.ApiToken, config.Telegram.ChatId, prices)
}
