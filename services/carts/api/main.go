package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/amane15/ecom_microservice/internal/platform"
	"github.com/amane15/ecom_microservice/services/carts/internal"
	"github.com/amane15/ecom_microservice/services/carts/internal/transport/http"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	cfg := internal.Config{
		Port:  4000,
		DbDSN: os.Getenv("DATABASE_URL"),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(cfg.DbDSN)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer db.Close()

	logger.Info("database connection pool established")

	app := internal.NewApplication(cfg, logger, db)
	handler := http.NewHandler(app)

	srv := platform.NewHTTPServer(":4003", handler.Routes())
	err = srv.Start()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
