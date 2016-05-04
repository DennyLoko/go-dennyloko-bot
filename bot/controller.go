package bot

import (
	"bytes"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/x-cray/logrus-prefixed-formatter"
)

// Controller is used to handle all the bot flows
type Controller struct {
	API *tgbotapi.BotAPI
	log *logrus.Logger
}

// NewController returns an new Controller object with the dependencies satisfied
func NewController(token string) (*Controller, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

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

	bot := &Controller{
		API: api,
		log: l,
	}

	l.Info("===================================")
	l.Infof("Bot ID: %d", bot.API.Self.ID)
	l.Infof("Bot name: %s", bot.API.Self.UserName)
	l.Info("===================================")

	return bot, nil
}

// Start the bot flow
func (b *Controller) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.API.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for u := range updates {
		b.log.Debug("Update received! Parsing...")
		b.log.Infof("[%s (%d)]: %s", u.Message.From.UserName, u.Message.From.ID, u.Message.Text)

		b.parseUpdate(u)
	}
}

func (b *Controller) parseUpdate(u tgbotapi.Update) {
	if u.Message.IsCommand() == true {
		switch u.Message.Command() {
		case "start":
			b.startCmd(u)
		case "help":
			b.helpCmd(u.Message)
		}
	} else {
		b.log.Debug("The human is trying to talk to me...")
		b.log.Debug("WHAT TO DO? WHAT TO DO?")
		b.log.Debug("Nothing.")

		m := tgbotapi.NewMessage(u.Message.Chat.ID, "\xE3\x80\xB0")
		b.API.Send(m)
	}
}

func (b *Controller) startCmd(u tgbotapi.Update) {
	msg := bytes.NewBufferString("")

	if u.Message.CommandArguments() == "start" {
		msg.WriteString("Welcome! I'm the DennyLoko's bot.\n")
		msg.WriteString("How can I help you?\n\n")
		msg.WriteString("You can type /help anytime to get a list of available commands.")
	} else {
		msg.WriteString("Welcome back! I missed you! \xF0\x9F\x98\x8D \xE2\x9D\xA4")
	}

	m := tgbotapi.NewMessage(u.Message.Chat.ID, msg.String())
	b.API.Send(m)
}

func (b *Controller) helpCmd(msg *tgbotapi.Message) {
	m := tgbotapi.NewMessage(msg.Chat.ID, "Sorry, there's no help topics yet... =(")
	b.API.Send(m)
}
