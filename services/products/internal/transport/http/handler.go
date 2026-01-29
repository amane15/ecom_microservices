package http

import (
	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/products/internal"
	"github.com/amane15/ecom_microservice/services/products/internal/service"
)

type Handler struct {
	app             *internal.Application
	httpErrRes      *httpx.ErrorResponse
	productService  *service.ProductService
	categoryService *service.CategoryService
}

func NewHandler(app *internal.Application, productService *service.ProductService, categoryService *service.CategoryService) *Handler {
	return &Handler{app: app, productService: productService, categoryService: categoryService}
}
