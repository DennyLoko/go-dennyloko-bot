package bot

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/x-cray/logrus-prefixed-formatter"
)

// Controller is used to handle all the bot flows
//
// List of supported commands by the bot
// help - Show the help topics
// currex - Do a currency conversion in the form: <value> <from> <to>
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
		case "currex", "c":
			b.currexCmd(u.Message)
		default:
			b.invalidCmd(u.Message)
		}
	} else {
		b.log.Debug("The human is trying to talk to me...")
		b.log.Debug("WHAT TO DO? WHAT TO DO?")
		b.log.Debug("Nothing.")

		m := tgbotapi.NewMessage(u.Message.Chat.ID, "\xE3\x80\xB0")
		b.API.Send(m)
	}
}

func (b *Controller) invalidCmd(msg *tgbotapi.Message) {
	m := tgbotapi.NewMessage(msg.Chat.ID, "I didn't understand this command, sorry.")
	b.API.Send(m)
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

func (b *Controller) currexCmd(m *tgbotapi.Message) {
	var msg, amount, from, to string
	var m2 tgbotapi.MessageConfig

	r := regexp.MustCompile("(\\d*(\\.?\\d*))(?:\\s)?(\\w{3})(\\s(?:to\\s)?(\\w{3}))?")
	q := r.FindStringSubmatch(m.CommandArguments())

	if len(q) == 0 || (len(q) == 6 && q[1] != " ") {
		msg = fmt.Sprint("Excuse me, but you've sent wrong parameters.\n")
		msg += fmt.Sprint("Please, try: `/currex amount from to`")

		m2 = tgbotapi.NewMessage(m.Chat.ID, msg)
		m2.ParseMode = "markdown"
		b.API.Send(m2)
		return
	}

	if len(q) == 4 || len(q) == 5 {
		amount = q[1]
		from = q[3]
		to = "BRL"
	} else if len(q) == 6 {
		amount = q[1]
		from = q[3]
		to = q[5]
	}

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	a, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}

	cx := &Currex{
		Amount: a,
		From:   from,
		To:     to,
		log:    b.log,
	}

	if err = cx.Validate(from); err != nil {
		m2 = tgbotapi.NewMessage(m.Chat.ID, err.Error())
		b.API.Send(m2)
		return
	}

	if err = cx.Validate(to); err != nil {
		m2 = tgbotapi.NewMessage(m.Chat.ID, err.Error())
		b.API.Send(m2)
		return
	}

	msg = fmt.Sprintf("Wait... I'm converting %.2f %s to %s", a, from, to)
	m2 = tgbotapi.NewMessage(m.Chat.ID, msg)
	b.API.Send(m2)

	s, f, t, err := cx.Convert()
	if err != nil {
		panic(err)
	}

	if s == true {
		msg = fmt.Sprintf("%s %.2f = %s %.2f", from, f, to, t)
	} else {
		msg = "I'm sorry, I wasn't able to do this conversion... =("
	}

	m2 = tgbotapi.NewMessage(m.Chat.ID, msg)
	b.API.Send(m2)
}
