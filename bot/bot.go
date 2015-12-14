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
	"strings"
)

type Bot struct {
	Tbot            *telebot.Bot
	Token           string
	Name            string
	Silent          bool
	LastMessage     *telebot.Message
	ChatContext     *chatContext
	Polling         time.Duration
	Wakeup          Action
	Sleep           Action
	CheckIfSleeping Action
	Random          RandomAction
	Actions         []Action `json:"actions"`
}

type chatContext struct {
	Sender     *telebot.User
	IsSleeping bool
}

type Action struct {
	Commands  []string
	Templates []string
}

type RandomAction struct {
	// RandomFrequency represents how frequently the bot will chat randomly. 1 means every msg, 10 means 1 out of 10 msgs and so on
	Frequency int
	Templates []string
}

func (b *Bot) Listen(subscription chan <- telebot.Message, timeout time.Duration) {
	b.Tbot.Listen(subscription, timeout)
}

func (b *Bot) SendMessage(recipient telebot.Recipient, message string, options *telebot.SendOptions) error {
	return b.Tbot.SendMessage(recipient, message, options)
}

func (b *Bot) SendAudio(recipient telebot.Recipient, audio *telebot.Audio, options *telebot.SendOptions) error {
	return b.Tbot.SendAudio(recipient, audio, options)
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
	b.ChatContext = &chatContext {Sender: &(b.LastMessage.Sender), IsSleeping:b.Silent}
}

func (b *Bot) check(commands []string) (bool) {

	txt := strings.ToLower(b.LastMessage.Text)
	for _, command := range commands {
		if strings.Contains(txt, command) {
			return true
		}
	}
	return false
}

func (b *Bot) timeToSaySomethingStupid(frequency int) (bool) {
	if frequency == 0 {
		return false
	}
	return rand.Intn(frequency) == 0
}

func (b* Bot) Start() {
	messages := make(chan telebot.Message)
	b.Listen(messages, b.Polling * time.Second)

	for message := range messages {
		b.SetLastMessage(message)

		if (b.check(b.Sleep.Commands)) {
			n := len(b.Sleep.Templates)
			b.SendMessage(message.Chat, b.SentenceFromTemplate(b.Sleep.Templates[rand.Intn(n)]), nil)
			b.Silent = true
		}

		if (b.check(b.CheckIfSleeping.Commands)) {
			n := len(b.CheckIfSleeping.Templates)
			b.SendMessage(message.Chat, b.SentenceFromTemplate(b.CheckIfSleeping.Templates[rand.Intn(n)]), nil)
		}

		if (!b.Silent) {
			for _, action := range b.Actions {
				if (b.check(action.Commands)) {
					n := len(action.Templates)
					answer := b.SentenceFromTemplate(action.Templates[rand.Intn(n)])
					if strings.HasSuffix(answer, ".ogg") {
						file, _ := telebot.NewFile(answer)
						audio := telebot.Audio{File: file}
						b.SendAudio(message.Chat, &audio, nil)
					} else {
						b.SendMessage(message.Chat, answer, nil)
					}
				}
			}
		}

		if (b.check(b.Wakeup.Commands)) {
			n := len(b.Wakeup.Templates)
			b.SendMessage(message.Chat, b.SentenceFromTemplate(b.Wakeup.Templates[rand.Intn(n)]), nil)
			b.Silent = false
		}

		if b.timeToSaySomethingStupid(b.Random.Frequency) {
			n := len(b.Random.Templates)
			b.SendMessage(message.Chat, b.SentenceFromTemplate(b.Random.Templates[rand.Intn(n)]), nil)
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

	bot.Tbot, err = telebot.NewBot(bot.Token)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Bot " + bot.Name + " started!")
	}

	return bot, err
}
