package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/carts/internal/data"
)

func (h *Handler) getCartHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	cart, err := h.app.Models.Carts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"cart": cart}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createCartHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID *int64 `json:"user_id"`
	}

	err := httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	if input.UserID == nil {
		h.httpErrRes.BadRequestResponse(w, r, errors.New("user_id must be provided"))
		return
	}

	cart, err := h.app.Models.Carts.Insert(*input.UserID)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"cart": cart}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
		return
	}
}

func (h *Handler) deleteCartHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	err = h.app.Models.Carts.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}
