package main

import (
	"encoding/json"
	"os"
	"github.com/tucnak/telebot"
	"math/rand"
	"time"
	"strings"
	"log"
	"text/template"
	"bytes"
	"github.com/ugol/botcreator/template/functions"
)

const dataFile = "data/bot/bot.json"
var ctx chatContext

type chatContext struct {
	Sender *telebot.User
	IsSleeping bool
}

type Bot struct {
	Token string
	Name  string
	Polling time.Duration
	Wakeup Action
	Sleep Action
	IsSleeping Action
	Actions []Action `json:"actions"`
}

type Action struct {
	Commands []string
	Templates  []string
}

func OpenBot() (*Bot, error) {

	rand.Seed(time.Now().UTC().UnixNano())
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var bot *Bot
	err = json.NewDecoder(file).Decode(&bot)

	return bot, err
}

func checkCommand(txt string, commands []string) (bool) {
	for _, command := range commands {
		if strings.Contains(txt, command) {
			return true
		}
	}
	return false
}

func sentenceFromTemplate(temp string) (string) {

	report, err := template.New("sentence").
	Funcs(functions.FunctionsMap()).
	Parse(temp)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	if err := report.Execute(&b, &ctx); err != nil {
		log.Fatal(err)
	}
	log.Println(b.String())
	return b.String()

}

func refreshContext(msg telebot.Message, isSleeping bool) {
	ctx = chatContext {Sender: &(msg.Sender), IsSleeping:isSleeping}
}

func main() {

	bot, err := OpenBot()

	if (err != nil) {
		log.Fatal(err)
	}

	tbot, err := telebot.NewBot(bot.Token)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Bot "+  bot.Name + " started!")
	}

	silent := false

	messages := make(chan telebot.Message)
	tbot.Listen(messages, bot.Polling * time.Second)

	for message := range messages {
		txt := strings.ToLower(message.Text)
		refreshContext(message, silent)

		if (checkCommand(txt, bot.Sleep.Commands)) {
			n := len(bot.Sleep.Templates)
			tbot.SendMessage(message.Chat, sentenceFromTemplate(bot.Sleep.Templates[rand.Intn(n)]), nil)
			silent = true
		}

		if (checkCommand(txt, bot.IsSleeping.Commands)) {
			n := len(bot.IsSleeping.Templates)
			tbot.SendMessage(message.Chat, sentenceFromTemplate(bot.IsSleeping.Templates[rand.Intn(n)]), nil)
		}

		if (!silent) {
			for _, action := range bot.Actions {
				if (checkCommand(txt, action.Commands)) {
					n := len(action.Templates)
					tbot.SendMessage(message.Chat, sentenceFromTemplate(action.Templates[rand.Intn(n)]), nil)
				}
			}
		}

		if (checkCommand(txt, bot.Wakeup.Commands)) {
			n := len(bot.Wakeup.Templates)
			tbot.SendMessage(message.Chat, sentenceFromTemplate(bot.Wakeup.Templates[rand.Intn(n)]), nil)
			silent = false
		}

	}

}