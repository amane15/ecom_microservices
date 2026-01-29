package http

import (
	"log/slog"

	"github.com/amane15/ecom_microservice/services/products/internal/service"
)

type Handler struct {
	productService  *service.ProductService
	categoryService *service.CategoryService
	logger          *slog.Logger
}

func NewHandler(productService *service.ProductService,
	categoryService *service.CategoryService,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		productService:  productService,
		categoryService: categoryService,
		logger:          logger,
	}
}
