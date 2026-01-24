package internal

import (
	"database/sql"
	"log/slog"

	"github.com/amane15/ecom_microservice/services/users/internal/data"
)

type Config struct {
	DbDSN string
	Port  int
}

type Application struct {
	Config
	Logger *slog.Logger
	Models data.Models
}

func NewApplication(cfg Config, logger *slog.Logger, db *sql.DB) *Application {
	return &Application{
		Logger: logger,
		Models: data.NewModels(db),
		Config: cfg,
	}
}
