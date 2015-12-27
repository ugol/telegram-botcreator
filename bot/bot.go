package bot

import (
	"encoding/json"
	"os"
	"github.com/tucnak/telebot"
	"math/rand"
	"time"
	"log"
	"bytes"
	"github.com/ugol/telegram-botcreator/template/functions"
	"text/template"
)

type Bot struct {
	*telebot.Bot
	Token           string
	Name            string
	Silent          bool
	LastMessage     *telebot.Message
	ChatContext     *chatContext
	Polling         time.Duration
	Wakeup          Action
	Sleep           Action
	CheckIfSleeping Action
	Actions         []Action `json:"actions"`
}

type chatContext struct {
	Sender     *telebot.User
	IsSleeping bool
}

func (b *Bot) SentenceFromTemplate(temp string) (string) {

	report, err := template.New("sentence").
	Funcs(functions.FunctionsMap()).
	Parse(temp)
	if err != nil {
		log.Fatal(err)
	}

	var bt bytes.Buffer
	if err := report.Execute(&bt, b.ChatContext); err != nil {
		log.Fatal(err)
	}
	log.Println(bt.String())
	return bt.String()

}

func (b *Bot) SetLastMessage(msg telebot.Message) {
	b.LastMessage = &msg
	b.ChatContext = &chatContext{Sender: &(b.LastMessage.Sender), IsSleeping:b.Silent}
}

func (b* Bot) Start() {
	messages := make(chan telebot.Message)
	b.Listen(messages, b.Polling * time.Second)

	for message := range messages {
		b.SetLastMessage(message)

		if (b.Sleep.canRun(b)) {
			b.Sleep.execute(b)
			b.Silent = true
		}

		if (b.CheckIfSleeping.canRun(b)) {
			b.CheckIfSleeping.execute(b)
		}

		if (!b.Silent) {
			for _, action := range b.Actions {
				if (action.canRun(b)) {
					action.execute(b)
				}
			}
		}

		if (b.Wakeup.canRun(b)) {
			b.Wakeup.execute(b)
			b.Silent = false
		}

	}
}

func OpenBot(name string) (*Bot, error) {

	rand.Seed(time.Now().UTC().UnixNano())
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var bot *Bot

	err = json.NewDecoder(file).Decode(&bot)
	if err != nil {
		return nil, err
	}

	bot.Bot, err = telebot.NewBot(bot.Token)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Bot " + bot.Name + " started!")
	}

	return bot, err
}
