package http

import (
	"log/slog"

	"github.com/amane15/ecom_microservice/services/users/internal/service"
)

type Handler struct {
	logger      *slog.Logger
	userService *service.UserService
}

func NewHandler(userService *service.UserService, logger *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}
