package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/amane15/ecom_microservice/internal/platform"
	"github.com/amane15/ecom_microservice/services/carts/internal/service"
	"github.com/amane15/ecom_microservice/services/carts/internal/store"
	"github.com/amane15/ecom_microservice/services/carts/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	dsn  string
	port string
}

func main() {
	godotenv.Load()

	cfg := config{
		port: os.Getenv("PORT"),
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

	cartStore := store.NewCartStore(dbpool)
	itemStore := store.NewItemStore(dbpool)
	cartService := service.NewCartService(cartStore, itemStore)
	handler := http.NewHandler(cartService, logger)

	addr := fmt.Sprintf(":%s", cfg.port)
	srv := platform.NewHTTPServer(addr, handler.Routes())
	logger.Info("starting server on", "addr", addr)
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
