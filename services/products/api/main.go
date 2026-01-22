package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	productpb "github.com/amane15/ecom_microservice/proto/product/v1"
	"github.com/amane15/ecom_microservice/services/proudcts/internal/data"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

type config struct {
	dsn  string
	port int
}

type application struct {
	config
	logger     *slog.Logger
	products   data.ProductModel
	variants   data.ProductVariantModel
	categories data.CategoryModel
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

	lis, err := net.Listen("tcp", "50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	productpb.RegisterProductServiceServer(grpcServer, &ProductGRPCServer{})
	logger.Info("Product grpc running on :50051")
	log.Fatal(grpcServer.Serve(lis))

	app := &application{
		config:     cfg,
		logger:     logger,
		products:   data.ProductModel{DB: db},
		variants:   data.ProductVariantModel{DB: db},
		categories: data.CategoryModel{DB: db},
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
