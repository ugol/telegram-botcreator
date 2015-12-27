package bot

import (
	"strings"
	"math/rand"
	"github.com/tucnak/telebot"
)

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

func (a *Action) execute(b *Bot) {

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

func (a *Action) canRun(b *Bot) (bool) {

	frequency := a.Frequency

	if frequency == -1 {
		return false
	}

	if frequency == 0 {
		// 0 is the default value and means every time, as 1
		frequency = 1
	}

	// rand.Intn(1) is always 0
	if rand.Intn(frequency) == 0 {

		if (a.Commands == nil) {
			return true
		}

		txt := strings.ToLower(b.LastMessage.Text)
		for _, command := range a.Commands {
			if strings.Contains(txt, command) {
				return true
			}
		}
	}
	return false
}

