package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	Configuration struct {
		NotificationService string `toml:"notification_service"`
		TradesFile          string `toml:"trades_file"`
		Binance             BinanceConfig
		Telegram            TelegramConfig
	}

	TelegramConfig struct {
		ApiToken string `toml:"api_token"`
		ChatId   int    `toml:"chat_id"`
	}

	BinanceConfig struct {
		Apikey    string `toml:"api_key"`
		Secretkey string `toml:"secret_key"`
	}
)

func ParseConfig(file string) Configuration {
	var data Configuration
	_, err := toml.DecodeFile(file, &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open TOML file")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return data
}
