package http

import "net/http"

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/products/{id}", h.getProductHandler)
	mux.HandleFunc("POST /v1/products", h.createProductHandler)
	mux.HandleFunc("PATCH /v1/products/{id}", h.updateProductHandler)
	mux.HandleFunc("PATCH /v1/products/{id}/mark", h.changeProductStatusHandler)
	mux.HandleFunc("PATCH /v1/products/{id}/default", h.setDefaultVariantHandler)
	mux.HandleFunc("GET /v1/products", h.listProductsHandler)
	mux.HandleFunc("GET /v1/products/{id}/variants", h.listProductVariantsHandler)

	mux.HandleFunc("GET /v1/variants/{id}", h.getVariantHandler)
	mux.HandleFunc("POST /v1/variants", h.createVariantHandler)
	mux.HandleFunc("PATCH /v1/variants/{id}", h.updateVariantHandler)
	mux.HandleFunc("DELETE /v1/variants/{id}", h.deleteVariantHandler)

	mux.HandleFunc("GET /v1/categories/{id}", h.getCategoryHandler)
	mux.HandleFunc("POST /v1/categories", h.createCategoryHandler)
	mux.HandleFunc("PATCH /v1/categories/{id}", h.updateCategoryHandler)
	mux.HandleFunc("PATCH /v1/categories/{id}/mark", h.markActiveHandler)
	mux.HandleFunc("GET /v1/categories", h.listCategoriesHandler)

	return mux
}
