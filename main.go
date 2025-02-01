package main

import (
	"fmt"
	govaluate "github.com/Knetic/govaluate"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var userInputs = make(map[int64]string)
var users = make(map[int64]string)

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Бот %s запущен", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			handleMessage(bot, update.Message)
		} else if update.CallbackQuery != nil {
			handleCallback(bot, update.CallbackQuery)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	switch message.Text {
	case "/start":
		users[chatID] = message.From.UserName
		msg := tgbotapi.NewMessage(chatID, "Привет! Я калькулятор-бот 🤖\nНажмите /calc, чтобы начать.")
		bot.Send(msg)
	case "/calc":
		userInputs[chatID] = "" // Очищаем ввод перед началом
		msg := tgbotapi.NewMessage(chatID, "Выберите числа и операцию:")
		msg.ReplyMarkup = getCalculatorKeyboard()
		bot.Send(msg)

	case "users</>":
		usersString := fmt.Sprintf("%v", users)
		msg := tgbotapi.NewMessage(chatID, usersString)
		bot.Send(msg)
	}
}

func handleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	chatID := callback.Message.Chat.ID
	data := callback.Data

	switch data {
	case "clear":
		userInputs[chatID] = ""
	case "=":
		result, err := evalExpression(userInputs[chatID])
		log.Print(result, err)
		if err != nil {
			userInputs[chatID] = "Ошибка"
		} else {
			userInputs[chatID] = result
		}
	default:
		userInputs[chatID] += data
	}

	msg := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, "Выражение: "+userInputs[chatID])
	kb := getCalculatorKeyboard()
	msg.ReplyMarkup = &kb
	bot.Send(msg)
}

func evalExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", fmt.Errorf("Ошибка: %v", err)
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", fmt.Errorf("Ошибка: %v", err)
	}

	return fmt.Sprintf("%v", result), nil
}

func getCalculatorKeyboard() tgbotapi.InlineKeyboardMarkup {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("7", "7"),
			tgbotapi.NewInlineKeyboardButtonData("8", "8"),
			tgbotapi.NewInlineKeyboardButtonData("9", "9"),
			tgbotapi.NewInlineKeyboardButtonData("/", "/"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("4", "4"),
			tgbotapi.NewInlineKeyboardButtonData("5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6", "6"),
			tgbotapi.NewInlineKeyboardButtonData("*", "*"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
			tgbotapi.NewInlineKeyboardButtonData("-", "-"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData(".", "."),
			tgbotapi.NewInlineKeyboardButtonData("=", "="),
			tgbotapi.NewInlineKeyboardButtonData("+", "+"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Очистить", "clear"),
		},
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
