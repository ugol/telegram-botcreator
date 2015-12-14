package main

import (
	"flag"
	"log"
	"github.com/ugol/telegram-botcreator/bot"
)

const dataFile = "data/bot/bot-test.json"

func main() {

	botJson := flag.String("json", dataFile, "JSON description of the BOT")
	flag.Parse()

	log.Println("Starting bot described in JSON file: [", *botJson, "]")
	bot, err := bot.OpenBot(*botJson)

	if (err != nil) {
		log.Fatal(err)
	} else {
		bot.Start()
	}

}