package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/carts/internal/data"
	"github.com/amane15/ecom_microservice/services/carts/internal/service"
)

func (h *Handler) getCartItemHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("get item request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	h.logger.Info("fetching a item", "id", id)
	item, err := h.cartService.GetItem(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error while fetching item", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}
	h.logger.Info("fetched item", "id", id)

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"item": item}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err, "id", id)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createCartItemHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("create item request received")
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	input := &service.CreateItemInput{}
	err = httpx.ReadJSON(w, r, input)
	if err != nil {
		h.logger.Error("error reading body", "error", err)
		httpx.BadRequestResponse(w, r, err)
		return
	}
	input.CartID = id

	h.logger.Info("creating an item")
	item, err := h.cartService.CreateItem(r.Context(), input)
	if err != nil {
		var valiationError *data.ValidationError
		switch {
		case errors.Is(err, data.ErrItemAlreadyExists):
			httpx.BadRequestResponse(w, r, errors.New("item already exists in a cart"))
		case errors.As(err, &valiationError):
			httpx.FailedValidationResponse(w, r, valiationError.Fields)
		default:
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}
	h.logger.Info("created an item")

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": item}, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err, "id", id)
		httpx.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		httpx.NotFoundResponse(w, r)
		return
	}

	err = h.cartService.DeleteItem(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			httpx.NotFoundResponse(w, r)
		default:
			h.logger.Error("error deleting an item", "error", err, "id", id)
			httpx.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		h.logger.Error("error while writing a json", "error", err, "id", id)
		httpx.ServerErrorResponse(w, r, err)
	}
}
