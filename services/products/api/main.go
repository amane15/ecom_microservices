package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/amane15/ecom_microservice/internal/platform"
	"github.com/amane15/ecom_microservice/services/products/internal"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
	"github.com/amane15/ecom_microservice/services/products/internal/store"
	"github.com/amane15/ecom_microservice/services/products/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := internal.Config{
		DbDSN: os.Getenv("DATABASE_URL"),
		Port:  4000,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := openDB(cfg.DbDSN)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer db.Close()

	logger.Info("database connection pool established")

	app := internal.NewApplication(cfg, logger)

	productStore := store.NewProductStore(db)
	variantStore := store.NewVariantStore(db)
	categoryStore := store.NewCategoryStore(db)

	productService := service.NewProductService(productStore, variantStore)
	categoryService := service.NewCategoryService(categoryStore)

	handler := http.NewHandler(app, productService, categoryService)

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
