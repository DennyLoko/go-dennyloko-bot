package main

import (
	"os"

	"github.com/DennyLoko/go-dennyloko-bot/bot"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	bot, _ := bot.NewController(os.Getenv("TELEGRAM_TOKEN"))
	bot.Start()
}
