# go-dennyloko-bot [WIP]
This is my personal bot which helps me on my daily tasks.
This is work is not done (and I think it'll never be), but you can use as you
wish.

It currently does nothing more than just connect to Telegram. New
functionalities will be added as needed.

## Installation
### Prerequisites
- go1.6+
- [govendor](https://github.com/kardianos/govendor)

### Installation
1. `git clone https://github.com/DennyLoko/go-dennyloko-bot.git`
3. `cd go-dennyloko-bot`
4. `govendor sync`
5. `go install`

## Running
As it connects to Telegram, you must have one bot token and inform it as
`TELEGRAM_TOKEN` environment variable. The bot can also read this variable from
a `.env` file.
