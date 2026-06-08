package main

import (
	"log/slog"
	"os"
	"web-chat/internal/logger"
)

func main() {
	cfg := logger.Config{
		Level:  logger.LevelFromString(os.Getenv("LOG_LEVEL")),
		Format: logger.FormatFromString(os.Getenv("LOG_FORMAT")),
	}

	log := logger.NewLogger(cfg)
	slog.SetDefault(log)

	slog.Info("server starting", "port", 8080)
}
