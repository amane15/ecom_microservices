package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/amane15/ecom_microservice/internal/platform"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
	"github.com/amane15/ecom_microservice/services/products/internal/store"
	"github.com/amane15/ecom_microservice/services/products/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type config struct {
	port int
	dsn  string
}

func main() {
	godotenv.Load()
	cfg := config{
		dsn:  os.Getenv("DATABASE_URL"),
		port: 4000,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dbpool, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer dbpool.Close()

	logger.Info("database connection pool established")

	productStore := store.NewProductStore(dbpool)
	variantStore := store.NewVariantStore(dbpool)
	categoryStore := store.NewCategoryStore(dbpool)

	productService := service.NewProductService(productStore, variantStore)
	categoryService := service.NewCategoryService(categoryStore)

	handler := http.NewHandler(productService, categoryService, logger)

	logger.Info("Starting server on address :4001")
	srv := platform.NewHTTPServer(":4001", handler.Routes())
	err = srv.Start()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer srv.Shutdown(context.Background())
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
