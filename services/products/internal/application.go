package internal

import (
	"log/slog"
)

type Config struct {
	DbDSN string
	Port  int
}

type Application struct {
	Config
	Logger *slog.Logger
}

func NewApplication(cfg Config, logger *slog.Logger) *Application {
	return &Application{
		Logger: logger,
		Config: cfg,
	}
}
