package main

import (
	"github.com/JxSam/max-contracts-bot/internal/app"
	"github.com/JxSam/max-contracts-bot/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	if err := app.Run(cfg); err != nil {
		panic(err)
	}
}
