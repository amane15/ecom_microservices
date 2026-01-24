package main

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/pkg/httpx"
	"github.com/amane15/ecom_microservice/pkg/validator"
	"github.com/amane15/ecom_microservice/services/carts/internal/data"
)

func (app *application) getCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}

	item, err := app.cartItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.httpErrRes.NotFoundResponse(w, r)
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusOK, httpx.Envelope{"item": item}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) createCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}

	var input struct {
		ProductID *int64 `json:"product_id"`
		VariantID *int64 `json:"variant_id"`
		Quantity  *int   `json:"quantity"`
	}

	err = httpx.ReadJSON(w, r, &input)
	if err != nil {
		app.httpErrRes.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.ProductID == nil && input.VariantID == nil && input.Quantity == nil {
		app.httpErrRes.BadRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v.Check(input.ProductID != nil, "product_id", "must be provided")
	v.Check(input.VariantID != nil, "variant_id", "variant_id must be provided")
	v.Check(input.Quantity != nil, "quantity", "must be provided")
	v.Check(*input.Quantity >= 1, "quantity", "must be greater that zero")

	if !v.Valid() {
		app.httpErrRes.FailedValidationResponse(w, r, v.Errors)
		return
	}

	item := &data.CartItem{
		CartID:    id,
		ProductID: *input.ProductID,
		VariantID: *input.VariantID,
		Quantity:  *input.Quantity,
	}

	err = app.cartItems.Insert(item)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrItemAlreadyExists):
			app.httpErrRes.BadRequestResponse(w, r, errors.New("item already exists in a cart"))
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusCreated, httpx.Envelope{"item": item}, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}

func (app *application) deleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.ReadIDParam(r)
	if err != nil {
		app.httpErrRes.NotFoundResponse(w, r)
		return
	}

	err = app.cartItems.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.httpErrRes.NotFoundResponse(w, r)
		default:
			app.httpErrRes.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = httpx.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.httpErrRes.ServerErrorResponse(w, r, err)
	}
}
