package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type (
	configuration struct {
		Notification_service string
		Trades_file          string
		Binance              binance_config
	}

	binance_config struct {
		Api_key    string
		Secret_key string
	}
)

func ParseConfig(file string) configuration {
	var data configuration
	_, err := toml.DecodeFile(file, &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return data
}
