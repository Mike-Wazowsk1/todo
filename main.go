package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var schedule map[int]int = make(map[int]int)

func test_fun(bot *tgbotapi.BotAPI, update tgbotapi.Update, period int, text string) {
	for {
		time.Sleep(time.Duration(period) * time.Second)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

func main() {
	schedule[0] = 60
	schedule[1] = 3600
	schedule[2] = 86400
	schedule[3] = 604800

	state := 0
	period := -1
	globalText := ""
	numRe := regexp.MustCompile("[0-9]+")

	bot, err := tgbotapi.NewBotAPI("5982995045:AAEUqBmxGiVuPPRawO67_i8COQYq3mPSsbI")
	if err != nil {
		panic(err)
	}

	updateConfig := tgbotapi.NewUpdate(0)

	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Создать привычку")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Редактировать привычку")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Мои привычки")),
	)
	var scheduleKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Минута")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Час")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("День")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Неделя")),
	)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "Создать привычку":
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Напиши текст привычки")
				state = 1

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		case "Минута":
			{
				period = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Через сколько минут?")
				state = 2

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}

			}
		case "Час":
			{
				period = 1
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Через сколько часов?")
				state = 2

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}

			}
		case "День":
			{
				period = 2
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Через сколько дней?")
				state = 2

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}

			}
		case "Неделя":
			{
				period = 3
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Через сколько недель?")
				state = 2

				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}

			}
		default:
			{
				if state == 1 && len(numRe.FindAllString(update.Message.Text, -1)) == 0 {
					text := "C какой периодичностью повторять"
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
					msg.ReplyMarkup = scheduleKeyboard
					globalText = update.Message.Text
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
					// go test_fun(bot, update, text,period)

				}
				if state == 2 {
					delta, err := strconv.ParseInt(update.Message.Text, 10, 64)
					if err != nil {
						panic(err)
					}
					go test_fun(bot, update, schedule[period]*int(delta), globalText)
				}
				if state == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
					msg.ReplyMarkup = numericKeyboard
					msg.ReplyToMessageID = update.Message.MessageID
					if _, err := bot.Send(msg); err != nil {
						panic(err)
					}
				}
			}
		}
		if len(numRe.FindAllString(update.Message.Text, -1)) > 1 {
			fmt.Println("nums")
			go test_fun(bot, update, period, globalText)
		}

	}
}
