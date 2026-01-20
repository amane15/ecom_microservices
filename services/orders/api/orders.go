package main

import (
	"errors"
	"net/http"

	"github.com/amane15/ecom_microservice/services/orders/internal/data"
)

func (app *application) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	order, err := app.orders.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch cart
	// Parse cart and cart_items into order and order orderItems
	// create full order
	// return response
}
func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {}
