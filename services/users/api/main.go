package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/amane15/ecom_microservice/internal/platform"
	"github.com/amane15/ecom_microservice/services/users/internal/platform/http"
	"github.com/amane15/ecom_microservice/services/users/internal/service"
	"github.com/amane15/ecom_microservice/services/users/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type config struct {
	dsn  string
	port int
}

func main() {
	godotenv.Load()

	cfg := config{
		port: 4000,
		dsn:  os.Getenv("DATABASE_URL"),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer db.Close()

	logger.Info("database connection pool established")

	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore)

	handler := http.NewHandler(userService, logger)

	logger.Info("Starting server on address :4000")
	srv := platform.NewHTTPServer(":4000", handler.Routes())

	err = srv.Start()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, nil
}
