package http

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/carts/internal/data"
)

func (h *Handler) getCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	item, err := h.app.Models.Items.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			h.httpErrRes.NotFoundResponse(w, r)
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"item": item}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) createCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	var input struct {
		ProductID *int64 `json:"product_id"`
		VariantID *int64 `json:"variant_id"`
		Quantity  *int32 `json:"quantity"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		h.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.ProductID == nil && input.VariantID == nil && input.Quantity == nil {
		h.httpErrRes.BadRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v.Check(input.ProductID != nil, "product_id", "must be provided")
	v.Check(input.VariantID != nil, "variant_id", "variant_id must be provided")
	v.Check(input.Quantity != nil, "quantity", "must be provided")
	v.Check(*input.Quantity >= 1, "quantity", "must be greater that zero")

	if !v.Valid() {
		h.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	item := &data.CartItem{
		CartID:    id,
		ProductID: sql.NullInt64{Int64: *input.ProductID, Valid: true},
		VariantID: sql.NullInt64{Int64: *input.VariantID, Valid: true},
		Quantity:  *input.Quantity,
	}

	err = h.app.Models.Items.Insert(item)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrItemAlreadyExists):
			h.httpErrRes.BadRequestResponse(w, r, errors.New("item already exists in a cart"))
		default:
			h.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": item}, nil)
	if err != nil {
		h.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		h.httpErrRes.NotFoundResponse(w, r)
		return
	}

	err = h.app.Models.Items.Delete(id)
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
