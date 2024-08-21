package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
)

func SendNotification(apitoken string, chatid int, prices []Price) {
	telegramService, err := telegram.New(apitoken)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialise Telegram service")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	telegramService.AddReceivers(int64(chatid))

	notify.UseServices(telegramService)

	dumped := []string{}
	for _, p := range prices {
		dumped = append(dumped, dump(p))
	}

	err = notify.Send(
		context.Background(),
		"Latest prices",
		strings.Join(dumped, "\n"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to send Telegram message")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
