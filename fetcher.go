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

	// This will fetch coins you use or have some amount in
	// acc, err := client.NewGetAccountService().Do(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // fmt.Println(acc)
	// // fmt.Print("\n\n")
	// // fmt.Println(acc.Balances)

	// assets := make([]string, 0)
	// for _, balance := range acc.Balances {
	// 	amount, err := strconv.ParseFloat(balance.Free, 64)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	if amount != 0 {
	// 		assets = append(assets, balance.Asset)
	// 	}
	// }
	// fmt.Println(assets)

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
