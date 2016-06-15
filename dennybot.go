package main

import (
	"os"

	"github.com/DennyLoko/go-dennyloko-bot/bot"
	"github.com/getsentry/raven-go"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	raven.SetDSN(os.Getenv("SENTRY_DSN"))
	raven.CapturePanic(func() {
		bot, _ := bot.NewController(os.Getenv("TELEGRAM_TOKEN"))
		bot.Start()
	}, nil)
}
