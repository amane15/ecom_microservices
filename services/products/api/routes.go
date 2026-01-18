package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/products/{id}", app.getProductHandler)
	mux.HandleFunc("POST /v1/products", app.createProductHandler)
	mux.HandleFunc("PATCH /v1/products/{id}", app.updateProductHandler)
	mux.HandleFunc("PATCH /v1/products/{id}/mark", app.changeProductStatusHandler)
	mux.HandleFunc("PATCH /v1/products/{id}/default", app.setDefaultVariantHandler)

	mux.HandleFunc("GET /v1/variants/{id}", app.getVariantHandler)
	mux.HandleFunc("POST /v1/variants", app.createVariantHandler)
	mux.HandleFunc("PATCH /v1/variants/{id}", app.updateVariantHandler)
	mux.HandleFunc("DELETE /v1/variants/{id}", app.deleteVariantHandler)

	mux.HandleFunc("GET /v1/categories/{id}", app.getCategoryHandler)
	mux.HandleFunc("POST /v1/categories", app.createCategoryHandler)
	mux.HandleFunc("PATCH /v1/categories/{id}", app.updateCategoryHandler)
	mux.HandleFunc("PATCH /v1/categories/{id}/mark", app.markActiveHandler)
	mux.HandleFunc("GET /v1/categories", app.listCategoriesHandler)

	return mux
}
