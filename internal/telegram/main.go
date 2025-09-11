package telegram

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	Token string
	Group int64
	Bot   *tgbotapi.BotAPI
}

func NewBot(token string, group int64) *Telegram {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	return &Telegram{
		Group: group,
		Bot:   bot,
	}
}

func (t *Telegram) SendVideo(title, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	fileBytes := tgbotapi.FileReader{
		Name:   file.Name(),
		Reader: file,
	}
	msg := tgbotapi.NewVideo(t.Group, fileBytes)
	msg.Caption = fmt.Sprintf("#%s - file: %s", title, file.Name())

	if _, err = t.Bot.Send(msg); err != nil {
		return fmt.Errorf("error sending video: %w", err)
	}

	return nil
}

func CheckEnvVars() error {
	tgToken := os.Getenv("TG_TOKEN")
	tgGroup := os.Getenv("TG_GROUP")

	if tgToken == "" || tgGroup == "" {
		return fmt.Errorf("missing environment variables")
	}

	return nil
}
