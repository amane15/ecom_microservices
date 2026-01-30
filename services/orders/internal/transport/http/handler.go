package http

import (
	"log/slog"

	"github.com/amane15/ecom_microservice/services/orders/internal/service"
)

type Handler struct {
	orderService *service.OrderService
	logger       *slog.Logger
}

func NewHandler(orderService *service.OrderService, logger *slog.Logger) *Handler {
	return &Handler{
		orderService: orderService,
		logger:       logger,
	}
}
