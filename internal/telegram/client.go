package telegram

import (
	"fmt"
	"log"
	"strconv"

	"silkroad/m/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Client struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramClient(cfg config.TelegramConfig) *Client {
	if cfg.BotToken == "" {
		logrus.Warn("Telegram bot token is empty, Telegram notifications will not work")
		return &Client{}
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Panic(err)
	}

	chatID, err := strconv.ParseInt(cfg.ChatID, 10, 64)
	if err != nil {
		logrus.Errorf("Error converting TELEGRAM_CHAT_ID: %v", err)
	}

	return &Client{
		bot:    bot,
		chatID: chatID,
	}
}

func (c *Client) SendTelegramMessage(message string) error {
	if c.bot == nil {
		logrus.Warn("Telegram bot is not configured, skipping message send")
		return nil
	}

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
