package utils

import (
	"log"
	"log/slog"
	"os"
)

func CreateLogger() *log.Logger {
	return slog.NewLogLogger(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}), slog.LevelDebug)
}
