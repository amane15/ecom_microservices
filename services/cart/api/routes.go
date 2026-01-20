package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/carts", app.createCartHandler)
	mux.HandleFunc("GET /v1/carts/{id}", app.getCartHandler)
	mux.HandleFunc("DELETE /v1/carts/{id}", app.deleteCartHandler)

	mux.HandleFunc("POST /v1/carts/{id}/items", app.createCartItemHandler)
	mux.HandleFunc("DELETE /v1/items/{id}", app.deleteCartItemHandler)

	return mux
}
