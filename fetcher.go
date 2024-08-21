package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

type (
	Price struct {
		Symbol           string
		FirstAmount      float64
		SecondAmount     float64
		PreviousExchange float64
		CurrentExchange  float64
		CurrentAmount    float64
		Difference       float64
	}
)

func GetData(apikey string, secretkey string, trades []Trade) []Price {
	client := binance.NewClient(apikey, secretkey)

	s := map[string]bool{}
	for _, trade := range trades {
		s[trade.Symbol] = true
	}

	symbols := make([]string, len(s))
	i := 0
	for k := range s {
		symbols[i] = k
		i++
	}

	prices, err := client.NewListPricesService().Symbols(symbols).Do(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pricesMap := map[string]float64{}
	for _, p := range prices {
		parsed, _ := strconv.ParseFloat(p.Price, 64)
		pricesMap[p.Symbol] = parsed
	}

	var calculatedTrades []Price
	for _, trade := range trades {
		currentExchange := pricesMap[trade.Symbol]
		previousExchange := trade.SecondAmount / trade.FirstAmount
		currentAmount := trade.FirstAmount * currentExchange
		difference := (currentExchange - previousExchange) / math.Abs(previousExchange) * 100

		calculatedTrade := Price{
			Symbol:           trade.Symbol,
			FirstAmount:      trade.FirstAmount,
			SecondAmount:     trade.SecondAmount,
			PreviousExchange: previousExchange,
			CurrentExchange:  currentExchange,
			CurrentAmount:    currentAmount,
			Difference:       difference,
		}

		calculatedTrades = append(calculatedTrades, calculatedTrade)
	}

	return calculatedTrades
}

func GetOpenTrades(apikey string, secretkey string, tradingpair string) []*binance.Order {
	client := binance.NewClient(apikey, secretkey)

	acc, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize Binance account service")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	assets := make([]string, 0)
	for _, balance := range acc.Balances {
		amount, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			continue
		}
		if amount != 0 {
			assets = append(assets, balance.Asset)
		}
	}

	allOrders := []*binance.Order{}
	for _, asset := range assets {
		if asset == "EUR" || asset == "USDC" {
			continue
		}

		symbol := fmt.Sprintf("%s%s", asset, tradingpair)
		orders, err := client.NewListOrdersService().Symbol(symbol).Do(context.Background())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to initialize Binance orders service")
			fmt.Fprintf(os.Stderr, "Symbol: %s\n", symbol)
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		allOrders = append(allOrders, orders...)
	}

	investedAmounts := map[string]float64{}
	for _, order := range allOrders {
		if order.Status != "FILLED" {
			fmt.Println("Status isn't filled, skipping")
			continue
		}

		invested, err := strconv.ParseFloat(order.CummulativeQuoteQuantity, 64)
		if err != nil {
			fmt.Println("Couldn't parse amount")
			continue
		}

		if order.Side == "BUY" {
			investedAmounts[order.Symbol] += invested
		} else if order.Side == "SELL" {
			investedAmounts[order.Symbol] -= invested
		} else {
			fmt.Fprintln(os.Stderr, "Unknown side")
			fmt.Fprintln(os.Stderr, order)
			os.Exit(1)
		}
	}

	filteredOrders := []*binance.Order{}
	for _, order := range allOrders {
		if investedAmounts[order.Symbol] <= 20 {
			continue
		}

		filteredOrders = append(filteredOrders, order)
	}

	return filteredOrders
}

func GetUsedCoins(apikey string, secretkey string) []string {
	client := binance.NewClient(apikey, secretkey)

	acc, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	assets := make([]string, 0)
	for _, balance := range acc.Balances {
		amount, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			continue
		}
		if amount != 0 {
			assets = append(assets, balance.Asset)
		}
	}
	fmt.Println(assets)

	return assets
}
