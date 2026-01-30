package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/amane15/ecom_microservice/internal/platform"
	"github.com/amane15/ecom_microservice/services/orders/internal/service"
	"github.com/amane15/ecom_microservice/services/orders/internal/store"
	"github.com/amane15/ecom_microservice/services/orders/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	port string
	dsn  string
}

func main() {
	godotenv.Load()

	cfg := config{
		port: "4000",
		dsn:  os.Getenv("DATABASE_URL"),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbpool, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer dbpool.Close()

	logger.Info("database connection pool established")

	orderStore := store.NewOrderStore(dbpool)
	itemStore := store.NewItemStore(dbpool)
	orderService := service.NewOrderService(orderStore, itemStore)
	handler := http.NewHandler(orderService, logger)

	logger.Info("starting server on address :4002")
	srv := platform.NewHTTPServer(":4002", handler.Routes())
	err = srv.Start()
	if err != nil {
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
