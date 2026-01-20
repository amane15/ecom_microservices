package main

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/services/carts/internal/data"
	"github.com/amane15/ecom_microservice/services/carts/internal/validator"
)

func (app *application) getCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	item, err := app.cartItems.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		ProductID *int64 `json:"product_id"`
		VariantID *int64 `json:"variant_id"`
		Quantity  *int   `json:"quantity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.ProductID == nil && input.VariantID == nil && input.Quantity == nil {
		app.badRequestResponse(w, r, errors.New("at least one field must be provided"))
		return
	}

	v.Check(input.ProductID != nil, "product_id", "must be provided")
	v.Check(input.VariantID != nil, "variant_id", "variant_id must be provided")
	v.Check(input.Quantity != nil, "quantity", "must be provided")
	v.Check(*input.Quantity >= 1, "quantity", "must be greater that zero")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
			app.badRequestResponse(w, r, errors.New("item already exists in a cart"))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCartItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.cartItems.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
