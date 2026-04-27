package utils

import (
	"log"
	"log/slog"
)

func CreateLogger() *log.Logger {
	return slog.NewLogLogger(&slog.TextHandler{}, slog.LevelDebug)
}
