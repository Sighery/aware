package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("Provide a configuration TOML file as the argument in CLI")
}

func main() {
	if len(os.Args) != 2 {
		help()
		return
	}

	arg := os.Args[1]

	if arg == "help" || arg == "-help" || arg == "--help" {
		help()
		return
	}

	config := ParseConfig(arg)
	trades := ParseTrades(config.Trades_file)

	fmt.Println(trades)

	prices := GetData(config.Binance.Api_key, config.Binance.Secret_key, trades)

	fmt.Printf("%+v\n", prices)
}
