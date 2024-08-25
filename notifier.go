package main

import (
	"context"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

func MeetsRules(rules []int, price Price) bool {
	for _, ruleInt := range rules {
		rule := float64(ruleInt)
		meetsRule := false
		if rule == 0 {
			meetsRule = true
		} else if rule < 0 {
			meetsRule = price.DifferencePerc < rule
		} else if rule > 0 {
			meetsRule = price.DifferencePerc > rule
		}

		if meetsRule {
			return true
		}
	}

	return false
}

func FilterTrades(rules []int, prices []Price, previous []Price) []Price {
	var mapping = map[string]map[string]Price{}

	for _, price := range previous {
		mapping[price.Id] = make(map[string]Price)
		mapping[price.Id]["Previous"] = price
	}

	for _, price := range prices {
		_, ok := mapping[price.Id]
		if !ok {
			mapping[price.Id] = make(map[string]Price)
		}
		mapping[price.Id]["Current"] = price
	}

	var result []Price
	for _, prices := range mapping {
		previous, previousExists := prices["Previous"]
		current, currentExists := prices["Current"]

		if !currentExists {
			continue
		}

		previousRules := false
		if !previousExists {
			previousRules = true
		} else {
			previousRules = MeetsRules(rules, previous)
		}

		currentRules := MeetsRules(rules, current)

		difference := math.Abs(previous.DifferencePerc - current.DifferencePerc)
		if !previousRules && currentRules {
			result = append(result, current)
		} else if previousRules && currentRules && difference >= 1 {
			result = append(result, current)
		}
	}

	return result

}

func SendMessage(apitoken string, chatid int, subject string, message string) {
	telegramService, err := telegram.New(apitoken)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialise Telegram service")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	telegramService.AddReceivers(int64(chatid))

	notify.UseServices(telegramService)

	err = notify.Send(
		context.Background(),
		subject,
		message,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to send Telegram message")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func SendNotification(apitoken string, chatid int, rules []int, prices []Price, previous []Price) {
	telegramService, err := telegram.New(apitoken)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialise Telegram service")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	telegramService.AddReceivers(int64(chatid))

	notify.UseServices(telegramService)

	var filtered []Price
	if previous == nil {
		for _, p := range prices {
			if MeetsRules(rules, p) {
				filtered = append(filtered, p)
			}
		}
	} else {
		result := FilterTrades(rules, prices, previous)
		if len(result) != 0 {
			filtered = result
		}
	}

	dumped := []string{}
	for _, trade := range filtered {
		dumped = append(dumped, dump(trade))
	}

	if len(dumped) == 0 {
		return
	}

	err = notify.Send(
		context.Background(),
		"Latest PnL",
		strings.Join(dumped, "\n"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to send Telegram message")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
