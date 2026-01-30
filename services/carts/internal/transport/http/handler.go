package http

import (
	"log/slog"

	"github.com/amane15/ecom_microservice/services/carts/internal/service"
)

type Handler struct {
	cartService *service.CartService
	logger      *slog.Logger
}

func NewHandler(cartService *service.CartService, logger *slog.Logger) *Handler {
	return &Handler{
		cartService: cartService,
		logger:      logger,
	}
}
