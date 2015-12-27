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

type Action struct {
	// Frequency represents how frequently the action will be executed.
	// -1 means never
	// 0 or 1 means every time
	// any number n > 1 means n%
	// example: 10 is a 10% frequency (1 out of 10)
	Frequency int
	Commands  []string
	Templates []string
}

func (a *Action) Execute(b *Bot) {

	n := len(a.Templates)
	answer := b.SentenceFromTemplate(a.Templates[rand.Intn(n)])
	if strings.HasSuffix(answer, ".ogg") {
		file, _ := telebot.NewFile(answer)
		audio := telebot.Audio{File: file}
		b.SendAudio(b.LastMessage.Chat, &audio, nil)
	} else {
		b.SendMessage(b.LastMessage.Chat, answer, nil)
	}

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

func (b *Bot) check(action Action) (bool) {

	frequency := action.Frequency

	if frequency == -1 {
		return false
	}

	if frequency == 0 {
		// 0 is the default value and means every time, as 1
		frequency = 1
	}

	// rand.Intn(1) is always 0
	if rand.Intn(frequency) == 0 {

		if (action.Commands == nil) {
			return true
		}

		txt := strings.ToLower(b.LastMessage.Text)
		for _, command := range action.Commands {
			if strings.Contains(txt, command) {
				return true
			}
		}
	}
	return false
}

func (b* Bot) Start() {
	messages := make(chan telebot.Message)
	b.Listen(messages, b.Polling * time.Second)

	for message := range messages {
		b.SetLastMessage(message)

		if (b.check(b.Sleep)) {
			b.Sleep.Execute(b)
			b.Silent = true
		}

		if (b.check(b.CheckIfSleeping)) {
			b.CheckIfSleeping.Execute(b)
		}

		if (!b.Silent) {
			for _, action := range b.Actions {
				if (b.check(action)) {
					action.Execute(b)
				}
			}
		}

		if (b.check(b.Wakeup)) {
			b.Wakeup.Execute(b)
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
