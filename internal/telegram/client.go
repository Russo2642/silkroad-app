package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strconv"
)

type Client struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramClient() *Client {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	chatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		logrus.Errorf("Error converting TELEGRAM_CHAT_ID: %v", err)
	}

	return &Client{
		bot:    bot,
		chatID: chatID,
	}
}

func (c *Client) SendTelegramMessage(message string) error {
	msg := tgbotapi.NewMessage(c.chatID, message)
	msg.ParseMode = "Markdown"
	_, err := c.bot.Send(msg)
	return err
}

func formatMessage(formID int, name, phone, email, description string) string {
	return fmt.Sprintf(
		"*ID формы*: %d\n*Имя*: %s\n*Телефон*: %s\n*Электронная почта*: %s\n*Текст*: %s",
		formID, name, phone, email, description,
	)
}
