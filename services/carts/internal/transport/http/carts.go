package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/carts/internal/data"
	"github.com/amane15/ecom_microservice/services/carts/internal/service"
)

func (h *Handler) getCartHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("get cart request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	h.logger.Info("fetching a cart", "id", id)
	cart, err := h.cartService.GetCart(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while fetching a cart", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}
	h.logger.Info("fetched a card", "id", id)

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"cart": cart}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createCartHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("create cart request received")
	input := &service.CreateCartInput{}

	err := httpx.ReadJSON(w, r, input)
	if err != nil {
		httpx.BadRequestResponse(w, r, err)
		return
	}

	h.logger.Info("creating a cart for user")
	cart, err := h.cartService.CreateCart(r.Context(), input)
	if err != nil {
		var validationError *data.ValidationError
		switch {
		case errors.As(err, &validationError):
			httpx.FailedValidationResponse(w, r, validationError.Fields)
		default:
			h.logger.Error("error while creating a cart", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}
	h.logger.Info("created a cart for a user", "user_id", cart.UserID)

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"cart": cart}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Handler) deleteCartHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("delete cart request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	h.logger.Info("deleting cart", "id", id)
	err = h.cartService.DeleteCart(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error deleting a cart", "error", err)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err)
		httpx.ServerErrorResponse(w, r, err)
	}
}
