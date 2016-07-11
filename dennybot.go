package main

import (
	"errors"
	"os"
	"time"

	"github.com/DennyLoko/go-dennyloko-bot/bot"
	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	client, _ := raven.New(os.Getenv("SENTRY_DSN"))

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	logl, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	l := logrus.New()
	l.Level = logl
	l.Formatter = new(prefixed.TextFormatter)

	defer func() {
		if rec := recover(); rec != nil {
			e := errors.New(rec.(string))

			interfaces := []raven.Interface{}
			exception := raven.NewException(e, raven.NewStacktrace(2, 3, nil))
			packet := raven.NewPacket(e.Error(), append(interfaces, exception)...)
			packet.Level = raven.FATAL

			go client.Capture(packet, nil)

			l.Errorf("PANIC: %s", e)
			time.Sleep(5 * time.Second)
		}
	}()

	bot, _ := bot.NewController(os.Getenv("TELEGRAM_TOKEN"), l)
	bot.Start()
}
