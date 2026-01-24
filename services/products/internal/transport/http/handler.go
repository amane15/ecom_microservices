package http

import (
	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/products/internal"
)

type Handler struct {
	app        *internal.Application
	httpErrRes *httpx.ErrorResponse
}

func NewHandler(app *internal.Application) *Handler {
	return &Handler{app: app}
}
