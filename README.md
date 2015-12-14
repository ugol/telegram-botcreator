# Telegram-botcreator
A nice chatter-bot creator for Telegram. You don't need to write code, just a simple JSON file. 

## Usage (from source)

```
go run main.go
```

the code will load the default JSON bot file data/bot/bot.json

Remember that you need to obtain your token from the [BotFather](https://telegram.me/BotFather) bot and add it to the JSON bot file.

## Usage (binary)

on Linux/OSX:

```
./telegram-botcreator
```

or, if you want to specify a different JSON bot file location than default (data/bot/bot.json)

```
./telegram-botcreator -json your-bot-file.json
```

on Windows:

```
telegram-botcreator.exe
```

or

```
telegram-botcreator.exe -json your-bot-file.json
```