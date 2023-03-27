package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var schedule map[int]int = make(map[int]int)

func test_fun(bot *tgbotapi.BotAPI, update tgbotapi.Update, period int, text string, stopCh chan struct{}, id int) {
	for {
		select {
		case <-stopCh:
			fmt.Printf("Stopping goroutine %d\n", id)
			return
		default:
			time.Sleep(time.Duration(period) * time.Second)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	stopChans := make(map[int]chan struct{})
	idToText := make(map[int]string)
	id := 0

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
		// tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Редактировать привычку")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Мои привычки")),
	)
	var scheduleKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Минута")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Час")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("День")),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Неделя")),
	)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		if update.CallbackQuery != nil {
			// Получаем ID привычки из инлайн-кнопки, на которую нажал пользователь
			if strings.Contains(update.CallbackQuery.Data, "habbit_") {
				var habbitMenu = tgbotapi.NewInlineKeyboardRow()
				// habbitMenu = append(habbitMenu, tgbotapi.NewInlineKeyboardButtonData("Назад", "back"))
				habbitMenu = append(habbitMenu, tgbotapi.NewInlineKeyboardButtonData("Отменить", "cancel_"+strings.TrimPrefix(update.CallbackQuery.Data, "habbit_")))
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Что будем делать?")

				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(habbitMenu)
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
			if strings.Contains(update.CallbackQuery.Data, "cancel_") {
				idStr := strings.TrimPrefix(update.CallbackQuery.Data, "cancel_")
				fmt.Println(idStr)

				idx, _ := strconv.Atoi(idStr)

				if err != nil {
					panic(err)
				}

				stopChan, found := stopChans[idx]
				if found {
					fmt.Println(idToText)

					close(stopChan)
					delete(stopChans, idx)
					delete(idToText, idx)
					fmt.Println(idToText)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Привычка отменена")
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}

			}

		}
		if update.Message != nil {

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
			case "Мои привычки":
				habbitMenu := tgbotapi.NewInlineKeyboardRow()
				for k, v := range idToText {
					habbitMenu = append(habbitMenu, tgbotapi.NewInlineKeyboardButtonData(v, "habbit_"+fmt.Sprint(k)))
					fmt.Println(v, k)
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Твои привычки")
				// habbitMenu = append(habbitMenu, tgbotapi.NewInlineKeyboardButtonData("Назад", "back"))
				// habbitMenu = append(habbitMenu, tgbotapi.NewInlineKeyboardButtonData("Отменить", "cancel"))
				if len(habbitMenu) > 0 {
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(habbitMenu)
				} else {
					msg.Text = "У тебя пока нет привычек"
				}
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}

			case "stop":

				stopChan, found := stopChans[id-1]
				if found {
					close(stopChan)
					delete(stopChans, id)
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
						stopChan := make(chan struct{})
						go test_fun(bot, update, schedule[period]*int(delta), globalText, stopChan, id)
						stopChans[id] = stopChan
						idToText[id] = globalText
						id++
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привычка создана!")
						msg.ReplyMarkup = numericKeyboard
						msg.ReplyToMessageID = update.Message.MessageID
						if _, err := bot.Send(msg); err != nil {
							panic(err)
						}
						state = 0
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
		}
	}

}
