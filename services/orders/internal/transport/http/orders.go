package http

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/services/orders/internal/data"
)

func (h *Handler) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	order, err := h.app.Models.Orders.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"order": order}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch cart
	// Parse cart and cart_items into order and order orderItems
	// create full order
	// return response
}
func (h *Handler) updateOrderHandler(w http.ResponseWriter, r *http.Request) {}
