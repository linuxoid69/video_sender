package telegram

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	kilobyte int64 = 1024
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

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error get stat for file: %w", err)
	}

	msg := tgbotapi.NewVideo(t.Group, fileBytes)
	msg.Caption = fmt.Sprintf("#%s - file: %s, size: %dKb", title, file.Name(), fileInfo.Size()*kilobyte*kilobyte)

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
