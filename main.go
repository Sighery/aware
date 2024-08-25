package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func help() {
	fmt.Println("Provide a configuration TOML file as the argument in CLI")
}

func dump(data interface{}) string {
	b, _ := json.MarshalIndent(data, "", "  ")
	return string(b)
}

func fetchTrades(config Configuration) []Trade {
	var trades []Trade
	if config.TradesFile != "" {
		trades = ParseTrades(config.TradesFile)
	} else {
		binanceTrades := GetOpenTrades(config.Binance.Apikey, config.Binance.Secretkey, config.TradingPair)
		trades = ConvertBinanceTrades(binanceTrades)
	}

	return trades
}

func formatPrices(config Configuration, prices []Price, previous []Price) string {
	var filtered []Price
	if previous == nil {
		for _, p := range prices {
			if MeetsRules(config.NotificationRules, p) {
				filtered = append(filtered, p)
			}
		}
	} else {
		result := FilterTrades(config.NotificationRules, prices, previous)
		if len(result) != 0 {
			filtered = result
		}
	}

	dumped := []string{}
	for _, trade := range filtered {
		dumped = append(dumped, dump(trade))
	}

	if len(dumped) == 0 {
		return ""
	}

	return strings.Join(dumped, "\n")
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

	pricesTicker := time.NewTicker(time.Duration(config.RefreshInterval) * time.Second)
	tradesTicker := time.NewTicker(time.Duration(config.TradesRefresh) * time.Second)

	trades := fetchTrades(config)
	var previous []Price

	go func() {
		for ; ; <-tradesTicker.C {
			fmt.Println("Updating trades...")
			trades = fetchTrades(config)
		}
	}()

	go func() {
		for ; ; <-pricesTicker.C {
			prices := GetData(config.Binance.Apikey, config.Binance.Secretkey, trades)

			fmt.Println("Latest PnL:")
			fmt.Println(dump(prices))

			message := formatPrices(config, prices, previous)

			if message != "" {
				SendMessage(config.Telegram.ApiToken, config.Telegram.ChatId, "", message)
			}
			previous = prices
		}
	}()

	select {}
}
