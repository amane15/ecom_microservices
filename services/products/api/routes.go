package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/products/{id}", app.getProductHandler)
	mux.HandleFunc("POST /v1/products", app.createProductHandler)
	mux.HandleFunc("PATCH /v1/products/{id}", app.updateProductHandler)

	return mux
}
