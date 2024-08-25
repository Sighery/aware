package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

type (
	Trade struct {
		Id           string
		Symbol       string
		FirstAmount  float64 `json:"first_amount"`
		SecondAmount float64 `json:"second_amount"`
		Remaining    float64 `json:"remaining"`
	}

	TradeHistory struct {
		bought float64
		order  *binance.Order
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

func ConvertBinanceTrades(orders []*binance.Order) []Trade {
	pairs := usedPairs(orders)

	trades := []Trade{}

	for _, pair := range pairs {
		pairOrders := gatherPair(orders, pair)

		holding := []TradeHistory{}
		for _, order := range pairOrders {
			bought, err := strconv.ParseFloat(order.ExecutedQuantity, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse bought")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if order.Side == "BUY" {
				holding = append(holding, TradeHistory{bought: bought, order: order})
			} else {
				holding = sellHolding(holding, bought)
			}
		}

		for _, trade := range holding {
			bought, err := strconv.ParseFloat(trade.order.ExecutedQuantity, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse bought")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			spent, err := strconv.ParseFloat(trade.order.CummulativeQuoteQuantity, 64)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to parse spent")
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			trades = append(trades, Trade{Id: strconv.FormatInt(trade.order.OrderID, 10), Symbol: trade.order.Symbol, FirstAmount: bought, SecondAmount: spent, Remaining: trade.bought})
		}
	}

	return trades
}

func sellHolding(holding []TradeHistory, sold float64) []TradeHistory {
	filtered := []TradeHistory{}
	matched := false
	// Try to cancel any matching trades first
	for _, trade := range holding {
		if sold == trade.bought {
			matched = true
			continue
		}
		filtered = append(filtered, trade)
	}

	if matched {
		return filtered
	}

	filtered = []TradeHistory{}

	currentSold := sold
	for _, trade := range holding {
		if currentSold <= trade.bought {
			trade.bought -= currentSold
			filtered = append(filtered, trade)
		} else if currentSold >= trade.bought {
			currentSold -= trade.bought
		}
	}

	return filtered
}

func usedPairs(orders []*binance.Order) []string {
	usedMap := map[string]bool{}
	for _, order := range orders {
		usedMap[order.Symbol] = true
	}

	used := []string{}
	for key := range usedMap {
		used = append(used, key)
	}

	return used
}

func gatherPair(orders []*binance.Order, symbol string) []*binance.Order {
	filtered := []*binance.Order{}
	for _, order := range orders {
		if order.Symbol != symbol {
			continue
		}
		filtered = append(filtered, order)
	}

	sort.Slice(filtered, func(i int, j int) bool {
		return filtered[i].Time < filtered[j].Time
	})

	return filtered
}
