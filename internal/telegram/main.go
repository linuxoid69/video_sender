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

func (t *Telegram) SendVideo(title, filePath string) (err error) {
	var file *os.File

	file, err = os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	fileBytes := tgbotapi.FileReader{
		Name:   file.Name(),
		Reader: file,
	}

	var fileInfo os.FileInfo

	fileInfo, err = file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get stat for file %s: %w", filePath, err)
	}

	msg := tgbotapi.NewVideo(t.Group, fileBytes)
	msg.Caption = fmt.Sprintf("#%s - file: %s, size: %dMb", title, file.Name(), fileInfo.Size()/kilobyte/kilobyte)

	if _, err = t.Bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send video file %s: %w", filePath, err)
	}

	return nil
}
