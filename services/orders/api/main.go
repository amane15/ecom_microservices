package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/orders/internal/data"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	dsn  string
	port int
}

type application struct {
	config
	logger     *slog.Logger
	httpErrRes *httpx.ErrorResponse
	orders     data.OrderModel
	orderItems data.OrderItemModel
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

	app := &application{
		config:     cfg,
		logger:     logger,
		httpErrRes: httpx.NewErrorResponseSender(logger),
		orders:     data.OrderModel{DB: db},
		orderItems: data.OrderItemModel{DB: db},
	}

	err = app.server()
	if err != nil {
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
