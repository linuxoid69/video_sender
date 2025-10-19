package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {
	Log := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(Log)
}
